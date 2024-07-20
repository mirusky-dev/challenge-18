package services

import (
	"context"
	"time"

	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/core/background/events"
	"github.com/gobp/gobp/models/dtos"
	"github.com/gobp/gobp/models/entities"
	"github.com/gobp/gobp/repositories"
	"github.com/hibiken/asynq"
	"golang.org/x/exp/slices"
)

type ITaskService interface {
	core.IBaseService[string, dtos.CreateTask, dtos.UpdateTask, entities.Task]
}

type taskService struct {
	taskRepository   repositories.ITaskRepository
	backgroundClient *asynq.Client
}

func NewTaskService(taskRepository repositories.ITaskRepository, backgroundClient *asynq.Client) ITaskService {
	return &taskService{
		taskRepository:   taskRepository,
		backgroundClient: backgroundClient,
	}
}

func (svc *taskService) Create(ctx context.Context, input dtos.CreateTask) (*entities.Task, *core.Exception) {
	appCtx, ok := core.FromContext(ctx)
	if !ok {
		return nil, core.MissingContext()
	}

	if slices.Contains(appCtx.Roles(), "manager") {
		return nil, core.BadRequest(core.WithMessage("managers can't create tasks"))
	}

	// checks if it's a valid input
	if err := input.Validate(ctx); err != nil {
		return nil, core.BadRequest(core.WithMessage(err.Error()))
	}

	task, err := svc.taskRepository.Create(ctx, entities.Task{
		Summary: input.Summary,
		UserID:  appCtx.UserID(),
	})
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (svc *taskService) GetByID(ctx context.Context, id string) (*entities.Task, *core.Exception) {

	task, err := svc.taskRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (svc *taskService) GetAll(ctx context.Context, limit, offset int) (*[]entities.Task, int, *core.Exception) {

	tasks, total, err := svc.taskRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return &tasks, total, nil
}

func (svc *taskService) Update(ctx context.Context, id string, input dtos.UpdateTask) (*entities.Task, *core.Exception) {
	appCtx, ok := core.FromContext(ctx)
	if !ok {
		return nil, core.MissingContext()
	}

	if slices.Contains(appCtx.Roles(), "manager") {
		return nil, core.BadRequest(core.WithMessage("managers can't update tasks"))
	}

	// checks if it's a valid input
	if err := input.Validate(ctx); err != nil {
		return nil, core.BadRequest(core.WithMessage(err.Error()))
	}

	task, err := svc.taskRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if task.UserID != appCtx.UserID() {
		return nil, core.Forbidden()
	}

	var changes entities.Task

	if input.Summary != nil {
		changes.Summary = *input.Summary
	}

	if task.PerformedAt == nil && input.Done != nil && *input.Done {
		now := time.Now()
		changes.PerformedAt = &now
	}

	task, err = svc.taskRepository.Update(ctx, id, changes)
	if err != nil {
		return nil, err
	}

	if input.Done != nil && *input.Done {
		task, _ := events.NewTaskCompleted(task.UserID, task.ID)
		svc.backgroundClient.Enqueue(task)
	}

	return &task, nil
}

func (svc *taskService) Delete(ctx context.Context, id string) *core.Exception {
	appCtx, ok := core.FromContext(ctx)
	if !ok {
		return core.MissingContext()
	}

	if !slices.Contains(appCtx.Roles(), "manager") {
		return core.BadRequest(core.WithMessage("tech can't delete tasks"))
	}

	err := svc.taskRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
