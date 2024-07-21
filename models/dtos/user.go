package dtos

import (
	"context"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mirusky-dev/challenge-18/models/entities"
)

type User struct {
	ID              string
	Username        string
	Email           string
	IsEmailVerified *bool
	Role            string
	ManagerID       *string

	CreatedAt *string
	UpdatedAt *string
	DeletedAt *string
}

func (u *User) FromEntity(e entities.User) {
	u.ID = e.ID
	u.Username = e.Username
	u.Email = e.Email
	u.IsEmailVerified = e.IsEmailVerified
	u.Role = e.Role
	u.ManagerID = e.ManagerID

	if !e.CreatedAt.IsZero() {
		formated := e.CreatedAt.Format(time.RFC3339)
		u.CreatedAt = &formated
	}

	if !e.UpdatedAt.IsZero() {
		formated := e.UpdatedAt.Format(time.RFC3339)
		u.UpdatedAt = &formated
	}

	if e.DeletedAt.Valid {
		formated := e.DeletedAt.Time.Format(time.RFC3339)
		u.CreatedAt = &formated
	}
}

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s CreateUser) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.Username, v.Required),
		v.Field(&s.Password, v.Required),
		v.Field(&s.Email, v.Required, is.Email),
	)
}

type UpdateUser struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

func (s UpdateUser) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.Email, v.Required.When(s.Username == nil)),
		v.Field(&s.Username, v.Required.When(s.Email == nil)),
	)
}
