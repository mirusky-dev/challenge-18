package dtos

import (
	"context"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type ChangePassword struct {
	Password string `json:"password"`
}

func (s ChangePassword) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.Password, v.Required),
	)
}
