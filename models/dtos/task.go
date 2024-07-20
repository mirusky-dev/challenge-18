package dtos

import (
	"context"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gobp/gobp/models/entities"
)

type Task struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
	UserID  string `json:"userId"`

	PerformedAt *string `json:"performedAt"`
	CreatedAt   *string `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt"`
	DeletedAt   *string `json:"deletedAt"`
}

func (t *Task) FromEntity(e entities.Task) {
	t.ID = e.ID
	t.Summary = e.Summary
	t.UserID = e.UserID
	if e.PerformedAt != nil {
		formated := e.PerformedAt.Format(time.RFC3339)
		t.PerformedAt = &formated
	}

	if !e.CreatedAt.IsZero() {
		formated := e.CreatedAt.Format(time.RFC3339)
		t.CreatedAt = &formated
	}

	if !e.UpdatedAt.IsZero() {
		formated := e.UpdatedAt.Format(time.RFC3339)
		t.UpdatedAt = &formated
	}

	if e.DeletedAt.Valid {
		formated := e.DeletedAt.Time.Format(time.RFC3339)
		t.CreatedAt = &formated
	}
}

type CreateTask struct {
	Summary string `json:"summary"`
}

func (s CreateTask) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.Summary, v.Required),
	)
}

type UpdateTask struct {
	Summary *string `json:"summary"`
	Done    *bool   `json:"done"`
}

func (s UpdateTask) Validate(ctx context.Context) error {
	return v.ValidateStruct(&s,
		v.Field(&s.Summary, v.Required.When(s.Done == nil)),
		v.Field(&s.Done, v.Required.When(s.Summary == nil)),
	)
}
