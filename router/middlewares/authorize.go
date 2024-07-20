package middlewares

import (
	"github.com/gobp/gobp/core"
	"github.com/gofiber/fiber/v2"
)

func Authorize(opts ...core.AuthorizationOption) fiber.Handler {

	return func(c *fiber.Ctx) error {

		appCtx, ok := core.FromContext(c.UserContext())
		if !ok {
			return core.MissingContext()
		}

		for _, opt := range opts {
			if err := opt.Evaluate(appCtx); err != nil {
				return err
			}
		}

		return c.Next()
	}
}
