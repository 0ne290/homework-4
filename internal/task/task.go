package task

import (
	"github.com/google/uuid"
	"homework-4/internal"
	"time"
)

type Status string

const (
	statusNew        Status = "new"
	statusInProgress Status = "in_progress"
	statusDone       Status = "done"
)

type Task struct {
	Uuid        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newTask(uuid uuid.UUID, title string, description *string, createdAt time.Time, updatedAt time.Time) *Task {
	return &Task{uuid, title, description, statusNew, createdAt, updatedAt}
}

func (task *Task) Update(updatedAt time.Time) error {
	switch task.Status {

	case statusNew:
		task.Status = statusInProgress
		task.UpdatedAt = updatedAt

		return nil

	case statusInProgress:
		task.Status = statusDone
		task.UpdatedAt = updatedAt

		return nil

	default:
		return &internal.InvariantViolationError{Message: "status is invalid"}
	}
}
