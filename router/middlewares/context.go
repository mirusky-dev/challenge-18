package middlewares

import (
	"github.com/gobp/gobp/core"
	"github.com/gofiber/fiber/v2"
)

// Context middleware creates a new userCtx for each request
func Context() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userCtx := core.NewUserCtx("", []string{}, []string{})

		c.SetUserContext(core.NewContext(c.Context(), userCtx))
		c.Set("Context-ID", userCtx.ID())

		return c.Next()
	}
}
