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

type IAuthService interface {
	Login(ctx context.Context, input dtos.Login) (string, string, time.Time, time.Time, *core.Exception)
	Logout(ctx context.Context, input dtos.Logout) *core.Exception
	RefreshToken(ctx context.Context, input dtos.RefreshToken) (string, string, time.Time, time.Time, *core.Exception)
	Register(ctx context.Context, input dtos.CreateUser) *core.Exception
	SendResetPassword(ctx context.Context, input dtos.SendResetPassword) *core.Exception
	VerifyResetPassword(ctx context.Context, input dtos.VerifyResetPassword) *core.Exception
}

type authService struct {
	mailer           mailer.Mailer
	passwordHasher   core.PasswordHasher
	userRepository   repositories.IUserRepository
	resetLinkStorage fiber.Storage
	tokenService     ITokenService // ????: Maybe it's not the best way
	config           env.Config
}

func NewAuthService(
	config env.Config,
	mailer mailer.Mailer,
	passwordHasher core.PasswordHasher,
	userRepository repositories.IUserRepository,
	resetLinkStorage fiber.Storage,
	tokenService ITokenService,
) IAuthService {
	return &authService{
		config:           config,
		mailer:           mailer,
		passwordHasher:   passwordHasher,
		userRepository:   userRepository,
		resetLinkStorage: resetLinkStorage,
		tokenService:     tokenService,
	}
}

func (svc *authService) Login(ctx context.Context, input dtos.Login) (string, string, time.Time, time.Time, *core.Exception) {
	if err := input.Validate(ctx); err != nil {
		return "", "", time.Time{}, time.Time{}, core.BadRequest(core.WithMessage(err.Error()))
	}

	user, errUser := svc.userRepository.FindByUsernameOrEmail(ctx, input.UsernameOrEmail, input.UsernameOrEmail)
	if errUser != nil {
		return "", "", time.Time{}, time.Time{}, core.Unauthorized(core.WithMessage("wrong username or password"))
	}

	ok, err := svc.passwordHasher.VerifyPassword(input.Password, user.Password)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, core.Unexpected(core.WithError(err))
	}

	if !ok {
		return "", "", time.Time{}, time.Time{}, core.Unauthorized(core.WithMessage("wrong username or password"))
	}

	token, refreshToken, expiresAt, refreshExpiresAt, exception := svc.tokenService.Issue(ctx, user.ID)
	if exception != nil {
		return "", "", time.Time{}, time.Time{}, exception
	}

	return token, refreshToken, expiresAt, refreshExpiresAt, nil
}

func (svc *authService) Logout(ctx context.Context, input dtos.Logout) *core.Exception {
	if exception := svc.tokenService.Revoke(ctx, input.TokenJTI, input.TokenExpiresAt); exception != nil {
		return exception
	}

	return nil
}

func (svc *authService) RefreshToken(ctx context.Context, input dtos.RefreshToken) (string, string, time.Time, time.Time, *core.Exception) {

	token, refreshToken, expiresAt, refreshExpiresAt, exception := svc.tokenService.Refresh(ctx, input.RefreshToken)
	if exception != nil {
		return "", "", time.Time{}, time.Time{}, exception
	}

	return token, refreshToken, expiresAt, refreshExpiresAt, nil
}

func (svc *authService) Register(ctx context.Context, input dtos.CreateUser) *core.Exception {
	if err := input.Validate(ctx); err != nil {
		return core.BadRequest(core.WithMessage(err.Error()))
	}

	found, errFind := svc.userRepository.FindByUsernameOrEmail(ctx, input.Username, input.Email)
	if errFind != nil || found.ID != "" {
		return core.BadRequest(core.WithMessage("username or email already taken"))
	}

	hash, err := svc.passwordHasher.HashPassword(input.Password)
	if err != nil {
		return core.Unexpected(core.WithError(err))
	}

	isVerified := true
	user, err := svc.userRepository.Create(ctx, entities.User{
		Username:        input.Username,
		Email:           input.Email,
		Password:        hash,
		IsEmailVerified: &isVerified,
	})
	if err != nil {
		return core.Unexpected(core.WithError(err))
	}

	// TODO: Use hermes
	template := `Dear %s,

WELCOME MESSAGE

Att. %s`
	welcomeMessage := fmt.Sprintf(template, user.Username, svc.config.EmailSenderName)
	svc.mailer.Send(mailer.Mail{
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
		Subject:   "Welcome to Rocket",
		PlainText: welcomeMessage,
	})

	return nil
}

func (svc *authService) SendResetPassword(ctx context.Context, input dtos.SendResetPassword) *core.Exception {

	user, err := svc.userRepository.FindByUsernameOrEmail(ctx, input.Email, input.Email)
	if err != nil {
		// Prevents email exposing
		return nil
	}

	id := uuid.New().String()
	resetLink := input.BaseURL + "/app/reset-password/" + id

	// TODO: Use Hermes
	template := "Dear %s,\n Password reset link: %s \nAtt. %s"
	resetPasswordMessage := fmt.Sprintf(template, user.Username, resetLink, svc.config.EmailSenderName)
	if err := svc.resetLinkStorage.Set(id, []byte(user.ID), 12*time.Hour); err != nil {
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
		Subject:   "Password Recovery",
		PlainText: resetPasswordMessage,
	}); err != nil {
		return core.Unexpected(core.WithError(err))
	}

	return nil
}

func (svc *authService) VerifyResetPassword(ctx context.Context, input dtos.VerifyResetPassword) *core.Exception {
	if err := input.Validate(ctx); err != nil {
		return core.BadRequest(core.WithMessage(err.Error()))
	}

	value, err := svc.resetLinkStorage.Get(input.ID)
	if err != nil {
		return core.Unexpected(core.WithError(err))
	}

	if value == nil {
		return core.BadRequest(core.WithMessage("reset code expired"))
	}

	// Prevents same code been used more than one time
	if err := svc.resetLinkStorage.Delete(input.ID); err != nil {
		return core.Unexpected(core.WithError(err))
	}

	userID := string(value)

	user, err := svc.userRepository.GetByID(ctx, userID)
	if err != nil || user.ID != userID {
		// user has been deleted or for some reason is inexistent
		return core.Unexpected()
	}

	hash, err := svc.passwordHasher.HashPassword(input.Password)
	if err != nil {
		return core.Unexpected(core.WithError(err))
	}

	err = svc.userRepository.ChangePassword(ctx, userID, hash, uuid.New().String())
	if err != nil {
		return core.Unexpected(core.WithError(err))
	}

	return nil
}
