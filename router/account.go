package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/mirusky-dev/challenge-18/core"
	"github.com/mirusky-dev/challenge-18/models/dtos"
)

func (ctrl Controller) changePassword(c *fiber.Ctx) error {
	var dto dtos.ChangePassword

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	exception := ctrl.accountService.ChangePassword(c.UserContext(), dto)
	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}

func (ctrl Controller) sendEmailVerificationLink(c *fiber.Ctx) error {
	exception := ctrl.accountService.SendVerificationEmail(c.UserContext(), c.BaseURL())

	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}

func (ctrl Controller) verifyEmailVerification(c *fiber.Ctx) error {
	id := c.Params("id")

	exception := ctrl.accountService.VerifyCode(c.UserContext(), id)
	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}
