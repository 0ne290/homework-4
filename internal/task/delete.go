package task

import (
	"github.com/gofiber/fiber/v2"
	"homework-4/internal/shared"
)

type DeleteResponse struct {
	Message string `json:"message"`
}

// @Summary Delete a task by ID
// @Produce json
// @Success 200 {object} DeleteResponse
// @Failure 400 {object} shared.Error400
// @Failure 500 {object} shared.Error500
// @Param id path int true "id"
// @Router /v1/tasks/{id} [delete]
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := c.Service.Delete(id)
	if err != nil {
		return err
	}

	return shared.Create200(ctx, response)
}

func (s *service) Delete(id int) (DeleteResponse, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return DeleteResponse{}, err
	}

	return DeleteResponse{"task deleted"}, nil
}

func (r *repository) Delete(id int) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	_, ok := r.tasks[id]
	if !ok {
		return &shared.InvariantViolationError{Message: "task not found"}
	}

	delete(r.tasks, id)

	return nil
}
