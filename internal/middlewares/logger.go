package middlewares

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"homework-4/internal"
)

func Logging(logger *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestUrl := "[" + ctx.Method() + "] " + ctx.OriginalURL()
		requestUuid := uuid.New().String()

		logger.Infow("begin", "url", requestUrl, "requestUuid", requestUuid, "requestBody", string(ctx.Request().Body()))

		ret := ctx.Next()
		if ret != nil {
			if errors.As(ret, &internal.NilOfInvariantViolationError) {
				ret = internal.Create400(ctx, &internal.Error400{Message: ret.Error()})

				logger.Infow("end", "requestUuid", requestUuid, "statusCode", 400, "responseBody", string(ctx.Response().Body()))
			} else {
				error500 := internal.NewError500(requestUrl, requestUuid)
				err := internal.Create500(ctx, &error500)

				logger.Errorw("end", "requestUuid", requestUuid, "statusCode", 500, "responseBody", string(ctx.Response().Body()), "error", ret)

				ret = err
			}
		} else {
			logger.Infow("end", "requestUuid", requestUuid, "statusCode", 200, "responseBody", string(ctx.Response().Body()))
		}

		return ret
	}
}
