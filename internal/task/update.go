package task

import (
	"github.com/gofiber/fiber/v2"
	"homework-4/internal/shared"
	"time"
)

type UpdateResponse struct {
	UpdatedAt time.Time `json:"updatedAt"`
	Status    Status    `json:"status"`
}

// @Summary Moves the task to the next status
// @Description Transitions of state machine: new -> in_progress -> done
// @Produce json
// @Success 200 {object} UpdateResponse
// @Failure 400 {object} shared.Error400
// @Failure 500 {object} shared.Error500
// @Param id path int true "id"
// @Router /v1/tasks/{id} [put]
func (c *Controller) Update(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := c.Service.Update(id)
	if err != nil {
		return err
	}

	return shared.Create200(ctx, response)
}

func (s *service) Update(id int) (UpdateResponse, error) {
	task, err := s.repository.Update(id)
	if err != nil {
		return UpdateResponse{}, err
	}

	return UpdateResponse{task.UpdatedAt, task.Status}, nil
}

func (r *repository) Update(id int) (Task, error) {
	task, ok := r.tasks[id]
	if !ok {
		return Task{}, &shared.InvariantViolationError{Message: "task not found"}
	}

	err := task.Update(time.Now())
	if err != nil {
		return Task{}, err
	}
	r.tasks[id] = task

	return task, nil
}
