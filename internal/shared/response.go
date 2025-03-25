package shared

import "github.com/gofiber/fiber/v2"

type Message struct {
	Message string `json:"message"`
}

type Error500 struct {
	Message     string `json:"message"`
	RequestUrl  string `json:"requestUrl"`
	RequestUuid string `json:"requestUuid"`
}

func NewError500(requestUrl, requestUuid string) Error500 {
	return Error500{"Please report the issue to technical support and attach this response body to your message.",
		requestUrl, requestUuid}
}

func Create200(ctx *fiber.Ctx, body any) error {
	return ctx.Status(fiber.StatusOK).JSON(body)
}

func Create400(ctx *fiber.Ctx, message Message) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(message)
}

func Create500(ctx *fiber.Ctx, err Error500) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(err)
}
