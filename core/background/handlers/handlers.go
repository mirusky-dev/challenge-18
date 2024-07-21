package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/mirusky-dev/challenge-18/core"
	"github.com/mirusky-dev/challenge-18/core/background/events"
	"github.com/mirusky-dev/challenge-18/core/env"
	"github.com/mirusky-dev/challenge-18/core/mailer"
	"github.com/mirusky-dev/challenge-18/models/dtos"
	"github.com/mirusky-dev/challenge-18/repositories"
)

type Controller struct {
	Config           env.Config
	Mailer           mailer.Mailer
	BackgroundClient *asynq.Client

	UserRepository repositories.IUserRepository
	TaskRepository repositories.ITaskRepository
}

func (ctrl *Controller) HandleTaskCompleted(ctx context.Context, t *asynq.Task) error {
	var p events.TaskCompletedPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	userEntity, err := ctrl.UserRepository.GetByID(ctx, p.UserID)
	if err != nil {
		return err
	}

	appCtx := core.NewUserCtx(userEntity.ID, []string{userEntity.Role}, []string{})
	taskEntity, err := ctrl.TaskRepository.GetByID(core.NewContext(ctx, appCtx), p.TaskID)
	if err != nil {
		return err
	}

	var manager, tech dtos.User
	manager.FromEntity(*userEntity.Manager)
	tech.FromEntity(userEntity)

	var task dtos.Task
	task.FromEntity(taskEntity)

	mail := mailer.Mail{
		From: mailer.EmailInfo{
			Name:  ctrl.Config.EmailSenderName,
			Email: ctrl.Config.EmailSender,
		},
		To: []mailer.EmailInfo{
			{
				Name:  manager.Username,
				Email: manager.Email,
			},
		},
		Subject: "Task Completed",
		PlainText: fmt.Sprintf("The tech %s (%s) performed the task %s on date %s",
			tech.Username,
			tech.ID,
			task.ID,
			*task.PerformedAt,
		),
	}

	return ctrl.Mailer.Send(mail)
}
