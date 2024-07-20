package events

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeTaskCompleted string = "task.completed"
)

type TaskCompletedPayload struct {
	UserID string
	TaskID string
}

func NewTaskCompleted(userID, taskID string) (*asynq.Task, error) {
	payload, _ := json.Marshal(TaskCompletedPayload{
		UserID: userID,
		TaskID: taskID,
	})
	return asynq.NewTask(TypeTaskCompleted, payload, asynq.TaskID(taskID), asynq.Unique(time.Hour)), nil
}
