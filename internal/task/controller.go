package task

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"homework-4/internal"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
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
		return &internal.InvariantViolationError{Message: "invalid request body format"}
	}

	response, err := c.service.Create(ctx.Context(), request)
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}

// @Summary Get all tasks
// @Accept json
// @Produce json
// @Success 200 {object} GetAllResponse
// @Failure 400 {object} shared.Error400
// @Failure 500 {object} shared.Error500
// @Router /v1/tasks [get]
func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	response, err := c.service.GetAll(ctx.Context())
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}

// @Summary Get task by ID
// @Produce json
// @Success 200 {object} GetByUuidResponse
// @Failure 400 {object} internal.Error400
// @Failure 500 {object} internal.Error500
// @Param uuid path string true "uuid" Format(uuid)
// @Router /v1/tasks/{uuid} [get]
func (c *Controller) GetByUuid(ctx *fiber.Ctx) error {
	taskUuid, err := uuid.Parse(ctx.Params("uuid"))
	if err != nil {
		return &internal.InvariantViolationError{Message: "invalid request body format"}
	}

	response, err := c.service.GetByUuid(ctx.Context(), taskUuid)
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}

// @Summary Delete a task by UUID
// @Produce json
// @Success 200 {object} DeleteResponse
// @Failure 400 {object} internal.Error400
// @Failure 500 {object} internal.Error500
// @Param uuid path string true "uuid" Format(uuid)
// @Router /v1/tasks/{uuid} [delete]
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	taskUuid, err := uuid.Parse(ctx.Params("uuid"))
	if err != nil {
		return &internal.InvariantViolationError{Message: "invalid request body format"}
	}

	response, err := c.service.Delete(ctx.Context(), taskUuid)
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}

// @Summary Moves the task to the next status
// @Description Transitions of state machine: new -> in_progress -> done
// @Produce json
// @Success 200 {object} UpdateResponse
// @Failure 400 {object} internal.Error400
// @Failure 500 {object} internal.Error500
// @Param id path string true "uuid" Format(uuid)
// @Router /v1/tasks/{uuid} [put]
func (c *Controller) Update(ctx *fiber.Ctx) error {
	taskUuid, err := uuid.Parse(ctx.Params("uuid"))
	if err != nil {
		return &internal.InvariantViolationError{Message: "invalid request body format"}
	}

	response, err := c.service.Update(ctx.Context(), taskUuid)
	if err != nil {
		return err
	}

	return internal.Create200(ctx, response)
}
