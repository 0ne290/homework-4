package middlewares

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"homework-4/internal/shared"
)

func Logging(logger *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestUrl := ctx.OriginalURL()
		requestUuid := uuid.New().String()

		logger.Infow("begin", "url", requestUrl, "requestUuid", requestUuid, "requestBody", string(ctx.Request().Body()))

		err := ctx.Next()
		if err != nil {
			if errors.As(err, &shared.NilOfInvariantViolationError) {
				_ = shared.Create400(ctx, &shared.Error400{Message: err.Error()})

				logger.Infow("end", "requestUuid", requestUuid, "statusCode", 400, "responseBody", string(ctx.Response().Body()))
			} else {
				error500 := shared.NewError500(requestUrl, requestUuid)
				_ = shared.Create500(ctx, &error500)

				logger.Errorw("end", "requestUuid", requestUuid, "statusCode", 500, "responseBody", string(ctx.Response().Body()))
			}
		} else {
			logger.Infow("end", "requestUuid", requestUuid, "statusCode", 200, "responseBody", string(ctx.Response().Body()))
		}

		return nil
	}
}
