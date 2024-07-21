package router

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/mirusky-dev/challenge-18/core"
	"github.com/mirusky-dev/challenge-18/models/dtos"
)

func (ctrl Controller) context(c *fiber.Ctx) error {
	userCtx, ok := core.FromContext(c.UserContext())

	if !ok {
		return core.MissingContext()
	}

	// Defining empty values to standardize front-end response
	// It should be updated by json.Marshal of userCtx
	result := fiber.Map{
		"roles":       []interface{}{},
		"permissions": []interface{}{},
	}

	token, ok := c.Locals("jwt").(*jwt.Token)
	if ok {
		claims := token.Claims.(*core.Claims)

		// Updates UserCtx values
		userCtx = userCtx.
			SetUserID(claims.Subject).
			SetRoles([]string{claims.Role})

	}

	b, _ := json.Marshal(userCtx)
	json.Unmarshal(b, &result)

	return c.JSON(result)
}

func (ctrl Controller) logout(c *fiber.Ctx) error {
	token := c.Locals("jwt").(*jwt.Token)
	claims := token.Claims.(*core.Claims)

	c.ClearCookie(refreshTokenKey)

	exception := ctrl.authService.Logout(c.UserContext(), dtos.Logout{
		TokenJTI:       claims.ID,
		TokenExpiresAt: claims.ExpiresAt.Time,
	})

	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}

func (ctrl Controller) register(c *fiber.Ctx) error {
	var dto dtos.CreateUser

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	exception := ctrl.authService.Register(c.UserContext(), dto)
	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}

func (ctrl Controller) login(c *fiber.Ctx) error {
	var dto dtos.Login

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	token, refreshToken, expiresAt, refreshExpiresAt, exception := ctrl.authService.Login(c.UserContext(), dto)
	if exception != nil {
		return exception
	}

	c.Cookie(&fiber.Cookie{
		Name:     refreshTokenKey,
		Value:    refreshToken,
		Expires:  refreshExpiresAt,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
	})

	return c.JSON(fiber.Map{
		"token":            token,
		"expiresAt":        expiresAt.Unix(),
		"refreshToken":     refreshToken,
		"refreshExpiresAt": refreshExpiresAt,
	})
}

func (ctrl Controller) refreshToken(c *fiber.Ctx) error {
	oldRefreshToken := c.Query("refresh_token", c.Cookies(refreshTokenKey))

	if oldRefreshToken == "" {
		return core.BadRequest(core.WithMessage("Missing refresh token"))
	}

	var dto = dtos.RefreshToken{
		RefreshToken: oldRefreshToken,
	}

	token, refreshToken, expiresAt, refreshExpiresAt, exception := ctrl.authService.RefreshToken(c.UserContext(), dto)

	if exception != nil {
		return exception
	}

	c.Cookie(&fiber.Cookie{
		Name:     refreshTokenKey,
		Value:    refreshToken,
		Expires:  refreshExpiresAt,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
	})

	return c.JSON(fiber.Map{
		"token":            token,
		"expiresAt":        expiresAt.Unix(),
		"refreshToken":     refreshToken,
		"refreshExpiresAt": refreshExpiresAt,
	})
}

func (ctrl Controller) sendResetPasswordLink(c *fiber.Ctx) error {
	var dto dtos.SendResetPassword

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	dto.BaseURL = c.BaseURL()

	exception := ctrl.authService.SendResetPassword(c.UserContext(), dto)

	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}

func (ctrl Controller) verifyResetPassword(c *fiber.Ctx) error {
	var dto dtos.VerifyResetPassword

	dto.ID = c.Params("id")

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	exception := ctrl.authService.VerifyResetPassword(c.UserContext(), dto)
	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}
