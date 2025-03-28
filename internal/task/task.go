package task

import (
	"homework-4/internal"
	"time"
)

type Status string

const (
	New        Status = "new"
	InProgress Status = "in_progress"
	Done       Status = "done"
)

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (task *Task) Update(updatedAt time.Time) error {
	switch task.Status {

	case New:
		task.Status = InProgress
		task.UpdatedAt = updatedAt

		return nil

	case InProgress:
		task.Status = Done
		task.UpdatedAt = updatedAt

		return nil

	default:
		return &internal.InvariantViolationError{Message: "status is invalid"}
	}
}
