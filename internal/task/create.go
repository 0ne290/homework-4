package task

import (
	"github.com/gofiber/fiber/v2"
	"homework-4/internal/shared"
	"time"
)

type CreateRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type CreateResponse struct {
	Id int `json:"id"`
}

// @Summary Create a new task with status "new"
// @Accept json
// @Produce json
// @Param createRequest body CreateRequest true "CreateRequest"
// @Success 200 {object} CreateResponse
// @Failure 400 {object} shared.Error400
// @Failure 500 {object} shared.Error500
// @Router /v1/tasks [post]
func (c *Controller) Create(ctx *fiber.Ctx) error {
	request := &CreateRequest{}
	if err := ctx.BodyParser(request); err != nil {
		return &shared.InvariantViolationError{Message: "invalid request body format"}
	}

	response, err := c.Service.Create(request)
	if err != nil {
		return err
	}

	return shared.Create200(ctx, response)
}

func (s *service) Create(request *CreateRequest) (CreateResponse, error) {
	task, err := s.repository.Create(request)
	if err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{task.Id}, nil
}

func (r *repository) Create(request *CreateRequest) (Task, error) {
	r.locker.Lock()

	task := Task{r.idCounter, request.Title, request.Description, New, time.Now(), time.Now()}
	r.tasks[r.idCounter] = task
	r.idCounter++

	r.locker.Unlock()

	return task, nil
}
