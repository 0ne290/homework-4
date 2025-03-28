package task

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"homework-4/internal"
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
	err := ctx.BodyParser(request)
	if err != nil {
		return err
	}

	response, err := c.Service.Create(request)
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}

func (s *Service) Create(ctx context.Context, request *CreateRequest) (CreateResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return CreateResponse{}, err
	}

	task, err := repository.Create(request)
	if err != nil {
		s.unitOfWork.Rollback(ctx, repository)

		return CreateResponse{}, err
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{task.Id}, nil
}

func (r *InMemoryRepository) Create(request *CreateRequest) (Task, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	task := Task{r.idCounter, request.Title, request.Description, New, time.Now(), time.Now()}
	r.tasks[r.idCounter] = task
	r.idCounter++

	return task, nil
}

func (r *PosgresRepository) Create(request *CreateRequest) (Task, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	task := Task{r.idCounter, request.Title, request.Description, New, time.Now(), time.Now()}
	r.tasks[r.idCounter] = task
	r.idCounter++

	return task, nil
}
