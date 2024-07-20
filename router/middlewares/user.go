package middlewares

import (
	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// User middleware resolves the current user based on JWT token
//
// also updated the userCtx values
func User(tokenService services.ITokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx, ok := core.FromContext(c.UserContext())
		if !ok {
			return core.MissingContext()
		}

		token := c.Locals("jwt").(*jwt.Token)
		claims := token.Claims.(*core.Claims)

		if exception := tokenService.IsRevoked(c.UserContext(), claims.ID); exception != nil {
			return exception
		}

		// Updates UserCtx values
		userCtx = userCtx.
			SetUserID(claims.Subject).
			SetRoles([]string{claims.Role})

		c.SetUserContext(core.NewContext(c.Context(), userCtx))

		return c.Next()
	}
}
