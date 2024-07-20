package dtos

import (
	"context"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Login struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

func (s Login) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.UsernameOrEmail, v.Required),
		v.Field(&s.Password, v.Required),
	)
}

type Logout struct {
	TokenJTI       string
	TokenExpiresAt time.Time
}

type VerifyResetPassword struct {
	ID       string `json:"-"`
	Password string `json:"password"`
}

func (s VerifyResetPassword) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.ID, v.Required),
		v.Field(&s.Password, v.Required),
	)
}

type RefreshToken struct {
	RefreshToken string
}

func (s RefreshToken) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.RefreshToken, v.Required),
	)
}

type SendResetPassword struct {
	BaseURL string `json:"-"`
	Email   string
}

func (s SendResetPassword) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.Email, is.Email, v.Required),
	)
}
