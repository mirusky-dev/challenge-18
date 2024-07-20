package services

import (
	"context"
	"fmt"
	"time"

	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/core/env"
	"github.com/gobp/gobp/core/mailer"
	"github.com/gobp/gobp/models/dtos"
	"github.com/gobp/gobp/models/entities"
	"github.com/gobp/gobp/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IAccountService interface {
	Me(ctx context.Context)
	ChangePassword(ctx context.Context, input dtos.ChangePassword) *core.Exception
	SendVerificationEmail(ctx context.Context, baseURL string) *core.Exception
	VerifyCode(ctx context.Context, id string) *core.Exception
}

type accountService struct {
	userRepository           repositories.IUserRepository
	passwordHasher           core.PasswordHasher
	emailVerificationStorage fiber.Storage
	mailer                   mailer.Mailer
	config                   env.Config
}

func NewAccountService(
	config env.Config,
	mailer mailer.Mailer,
	passwordHasher core.PasswordHasher,
	emailVerificationStorage fiber.Storage,
	userRepository repositories.IUserRepository,
) IAccountService {
	return &accountService{
		userRepository:           userRepository,
		passwordHasher:           passwordHasher,
		emailVerificationStorage: emailVerificationStorage,
		mailer:                   mailer,
		config:                   config,
	}
}

func (svc *accountService) Me(ctx context.Context) {}

func (svc *accountService) ChangePassword(ctx context.Context, input dtos.ChangePassword) *core.Exception {
	appCtx, ok := core.FromContext(ctx)
	if !ok {
		return core.MissingContext()
	}

	if _, err := svc.userRepository.GetByID(ctx, appCtx.UserID()); err != nil {
		return core.Unexpected(core.WithError(err))
	}

	hash, err := svc.passwordHasher.HashPassword(input.Password)
	if err != nil {
		return core.Unexpected()
	}

	if err := svc.userRepository.ChangePassword(ctx, appCtx.UserID(), hash, uuid.New().String()); err != nil {
		return core.Unexpected()
	}

	return nil
}

func (svc *accountService) SendVerificationEmail(ctx context.Context, baseURL string) *core.Exception {
	appCtx, _ := core.FromContext(ctx)

	user, err := svc.userRepository.GetByID(ctx, appCtx.UserID())
	if err != nil {
		// user has been deleted or for some reason is inexistent
		return nil
	}

	id := uuid.New().String()
	emailVerificationLink := baseURL + "/app/verify-email/" + id

	// TODO: Use Hermes
	template := `Dear %s,

Email verification link: %s

Att. %s`

	emailVerificationMessage := fmt.Sprintf(
		template,
		user.Username,
		emailVerificationLink,
		svc.config.EmailSenderName,
	)
	if err := svc.emailVerificationStorage.Set(id, []byte(user.ID), 12*time.Hour); err != nil {
		return core.Unexpected(core.WithError(err))
	}

	if err := svc.mailer.Send(mailer.Mail{
		From: mailer.EmailInfo{
			Name:  svc.config.EmailSenderName,
			Email: svc.config.EmailSender,
		},
		To: []mailer.EmailInfo{
			{
				Name:  user.Username,
				Email: user.Email,
			},
		},
		Subject:   "Email Verification",
		PlainText: emailVerificationMessage,
	}); err != nil {
		return core.Unexpected(core.WithError(err))
	}

	return nil
}

func (svc *accountService) VerifyCode(ctx context.Context, id string) *core.Exception {
	value, err := svc.emailVerificationStorage.Get(id)
	if err != nil {
		return core.Unexpected(core.WithError(err))
	}

	if value == nil {
		return core.BadRequest(core.WithMessage("verification code expired"))
	}

	// Prevents same code been used more than one time
	if err := svc.emailVerificationStorage.Delete(id); err != nil {
		return core.Unexpected(core.WithError(err))
	}

	userID := string(value)

	user, _ := svc.userRepository.GetByID(ctx, userID)

	if user.ID != userID {
		// user has been deleted or for some reason is inexistent
		return core.Unexpected()
	}

	isEmailVerified := true
	if _, err := svc.userRepository.Update(ctx, userID, entities.User{IsEmailVerified: &isEmailVerified}); err != nil {
		return core.Unexpected(core.WithError(err))
	}

	// TODO: congrats for having email verified

	return nil
}
