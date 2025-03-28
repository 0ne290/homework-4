package task

import (
	"github.com/gofiber/fiber/v2"
	"homework-4/internal"
)

type GetAllResponse struct {
	Tasks []Task `json:"tasks"`
}

// @Summary Get all tasks
// @Accept json
// @Produce json
// @Success 200 {object} GetAllResponse
// @Failure 400 {object} shared.Error400
// @Failure 500 {object} shared.Error500
// @Router /v1/tasks [get]
func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	response, err := c.Service.GetAll()
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}

func (s *Service) GetAll() (GetAllResponse, error) {
	tasks, err := s.repository.GetAll()
	if err != nil {
		return GetAllResponse{}, err
	}

	return GetAllResponse{tasks}, nil
}

func (r *InMemoryRepository) GetAll() ([]Task, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	ret := make([]Task, 0, len(r.tasks))

	for _, task := range r.tasks {
		ret = append(ret, task)
	}

	return ret, nil
}

type GetByIdResponse struct {
	Task Task `json:"task"`
}

// @Summary Get task by ID
// @Produce json
// @Success 200 {object} GetByIdResponse
// @Failure 400 {object} shared.Error400
// @Failure 500 {object} shared.Error500
// @Param id path int true "id"
// @Router /v1/tasks/{id} [get]
func (c *Controller) GetById(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	response, err := c.Service.GetById(id)
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}

func (s *Service) GetById(id int) (GetByIdResponse, error) {
	task, err := s.repository.GetById(id)
	if err != nil {
		return GetByIdResponse{}, err
	}

	return GetByIdResponse{task}, nil
}

func (r *InMemoryRepository) GetById(id int) (Task, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return Task{}, &internal.InvariantViolationError{Message: "task not found"}
	}

	return task, nil
}
