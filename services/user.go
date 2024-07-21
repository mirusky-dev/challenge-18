package services

import (
	"context"

	"github.com/mirusky-dev/challenge-18/core"
	"github.com/mirusky-dev/challenge-18/models/dtos"
	"github.com/mirusky-dev/challenge-18/models/entities"
	"github.com/mirusky-dev/challenge-18/repositories"
)

type IUserService interface {
	core.IBaseService[string, dtos.CreateUser, dtos.UpdateUser, entities.User]
}

type userService struct {
	userRepository repositories.IUserRepository
	passwordHasher core.PasswordHasher
}

func NewUserService(userRepository repositories.IUserRepository, passwordHasher core.PasswordHasher) IUserService {
	return &userService{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (svc *userService) Create(ctx context.Context, input dtos.CreateUser) (*entities.User, *core.Exception) {

	// checks if it's a valid input
	if err := input.Validate(ctx); err != nil {
		return nil, core.BadRequest(core.WithMessage(err.Error()))
	}

	// run the creation flow
	found, _ := svc.userRepository.FindByUsernameOrEmail(ctx, input.Username, input.Email)
	if found.ID != "" {
		return nil, core.BadRequest(core.WithMessage("username or email already taken"))
	}

	hashedPassword, err := svc.passwordHasher.HashPassword(input.Password)
	if err != nil {
		return nil, core.Unexpected()
	}

	user, errUser := svc.userRepository.Create(ctx, entities.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	})
	if errUser != nil {
		return nil, errUser
	}

	return &user, nil
}

func (svc *userService) GetByID(ctx context.Context, id string) (*entities.User, *core.Exception) {

	user, err := svc.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (svc *userService) GetAll(ctx context.Context, limit, offset int) (*[]entities.User, int, *core.Exception) {

	users, total, err := svc.userRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return &users, total, nil
}

func (svc *userService) Update(ctx context.Context, id string, input dtos.UpdateUser) (*entities.User, *core.Exception) {

	// checks if it's a valid input
	if err := input.Validate(ctx); err != nil {
		return nil, core.BadRequest(core.WithMessage(err.Error()))
	}

	user, err := svc.userRepository.GetByID(ctx, id)
	if err != nil || user.ID == "" {
		return nil, core.NotFound()
	}

	var username string
	var email string

	if input.Username == nil {
		input.Username = &username
	}

	if input.Email == nil {
		input.Email = &email
	}

	user, _ = svc.userRepository.FindByUsernameOrEmail(ctx, *input.Username, *input.Email)
	if user.ID != id {
		return nil, core.BadRequest(core.WithMessage("username or email already taken"))
	}

	user, err = svc.userRepository.Update(ctx, id, entities.User{Username: *input.Username, Email: *input.Email})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (svc *userService) Delete(ctx context.Context, id string) *core.Exception {

	err := svc.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
