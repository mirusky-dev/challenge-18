package router

import (
	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/models/dtos"
	"github.com/gofiber/fiber/v2"
)

func (ctrl Controller) createTask(c *fiber.Ctx) error {
	var input dtos.CreateTask

	if err := c.BodyParser(&input); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	task, exception := ctrl.taskService.Create(c.UserContext(), input)
	if exception != nil {
		return exception
	}

	var result dtos.Task
	result.FromEntity(*task)

	return c.JSON(result)
}

func (ctrl Controller) getTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")

	task, exception := ctrl.taskService.GetByID(c.UserContext(), id)
	if exception != nil {
		return exception
	}

	if task == nil {
		return core.NotFound()
	}

	var result dtos.Task
	result.FromEntity(*task)

	return c.JSON(result)
}

func (ctrl Controller) getAllTasks(c *fiber.Ctx) error {
	var input core.PaginationParams

	if err := c.QueryParser(&input); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	input.Default()

	tasks, total, exception := ctrl.taskService.GetAll(c.UserContext(), *input.Limit, *input.Offset)
	if exception != nil {
		return exception
	}

	var items []dtos.Task
	for _, task := range *tasks {
		var result dtos.Task
		result.FromEntity(task)
		items = append(items, result)
	}

	response := core.Page(items, int(total), *input.Limit, *input.Offset)
	return c.JSON(response)
}

func (ctrl Controller) updateTask(c *fiber.Ctx) error {
	var dto dtos.UpdateTask

	if err := c.BodyParser(&dto); err != nil {
		return core.BadRequest(core.WithError(err))
	}

	id := c.Params("id")

	task, exception := ctrl.taskService.Update(c.UserContext(), id, dto)
	if exception != nil {
		return exception
	}

	var result dtos.Task
	result.FromEntity(*task)

	return c.JSON(result)
}

func (ctrl Controller) deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	exception := ctrl.taskService.Delete(c.UserContext(), id)
	if exception != nil {
		return exception
	}

	return c.SendStatus(200)
}
