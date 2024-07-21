package router

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mirusky-dev/challenge-18/core"
)

func (ctrl Controller) empty(c *fiber.Ctx) error {
	return nil
}

func (ctrl Controller) pong(message string) fiber.Handler {
	return func(c *fiber.Ctx) error { return c.SendString(message) }
}

func (ctrl Controller) friendlyError(c *fiber.Ctx) error {
	var q core.Exception

	if err := c.QueryParser(&q); err != nil {
		return core.BadRequest()
	}

	var opts []core.UserFriendlyExceptionOption

	if q.Status != 0 {
		opts = append(opts, core.WithStatus(q.Status))
	}

	if q.Code != "" {
		opts = append(opts, core.WithCode(q.Code))
	}

	if q.Message != "" {
		opts = append(opts, core.WithMessage(q.Message))
	}

	if q.Err != "" {
		opts = append(opts, core.WithError(errors.New(q.Err)))
	}

	if q.Severity != "" {
		opts = append(opts, core.WithSeverity(q.Severity))
	}

	return core.UserFriendlyException(opts...)
}
