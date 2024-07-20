package router

import (
	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/models/dtos"
	"github.com/gofiber/fiber/v2"
)

func (ctrl Controller) createUser(c *fiber.Ctx) error {
	var dto dtos.CreateUser

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	user, exception := ctrl.userService.Create(c.UserContext(), dto)

	if exception != nil {
		return exception
	}

	return c.JSON(fiber.Map{
		"id":              user.ID,
		"username":        user.Username,
		"email":           user.Email,
		"role":            user.Role,
		"isEmailVerified": *user.IsEmailVerified,
		"createdAt":       user.CreatedAt,
		"updatedAt":       user.UpdatedAt,
		"deletedAt":       user.DeletedAt,
	})
}

func (ctrl Controller) getUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, exception := ctrl.userService.GetByID(c.UserContext(), id)
	if exception != nil {
		return exception
	}

	if user == nil {
		return core.NotFound()
	}

	return c.JSON(fiber.Map{
		"id":              user.ID,
		"username":        user.Username,
		"email":           user.Email,
		"role":            user.Role,
		"isEmailVerified": *user.IsEmailVerified,
		"createdAt":       user.CreatedAt,
		"updatedAt":       user.UpdatedAt,
		"deletedAt":       user.DeletedAt,
	})
}

func (ctrl Controller) getAllUsers(c *fiber.Ctx) error {
	var input core.PaginationParams

	if err := c.QueryParser(&input); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	input.Default()

	users, total, exception := ctrl.userService.GetAll(c.UserContext(), *input.Limit, *input.Offset)
	if exception != nil {
		return exception
	}

	var items []fiber.Map
	for _, user := range *users {
		items = append(items, fiber.Map{
			"id":              user.ID,
			"username":        user.Username,
			"email":           user.Email,
			"role":            user.Role,
			"isEmailVerified": *user.IsEmailVerified,
			"createdAt":       user.CreatedAt,
			"updatedAt":       user.UpdatedAt,
			"deletedAt":       user.DeletedAt,
		})
	}

	response := core.Page(items, int(total), *input.Limit, *input.Offset)
	return c.JSON(response)
}

func (ctrl Controller) updateUser(c *fiber.Ctx) error {
	var dto dtos.UpdateUser

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	id := c.Params("id")

	user, exception := ctrl.userService.Update(c.UserContext(), id, dto)

	if exception != nil {
		return exception
	}

	return c.JSON(fiber.Map{
		"id":              user.ID,
		"username":        user.Username,
		"email":           user.Email,
		"role":            user.Role,
		"isEmailVerified": *user.IsEmailVerified,
		"createdAt":       user.CreatedAt,
		"updatedAt":       user.UpdatedAt,
		"deletedAt":       user.DeletedAt,
	})
}

func (ctrl Controller) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	exception := ctrl.userService.Delete(c.UserContext(), id)
	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}
