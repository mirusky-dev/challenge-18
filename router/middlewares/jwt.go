package middlewares

import (
	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/core/env"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// JWT Resolves jwt
func JWT(config env.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.JWTSecret),
		ContextKey: "jwt",
		Claims:     &core.Claims{},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return core.BadRequest(core.WithMessage(err.Error()))
			}
			return core.Unauthorized(core.WithMessage("Invalid or expired JWT"))
		},
	})
}
