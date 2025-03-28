package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	_ "homework-4/docs"
	"homework-4/internal/shared"
	"homework-4/internal/shared/middlewares"
	"homework-4/internal/task"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title Homework 4 API
// @version 1.0
// @description Task CRUD
// @host localhost:8080
// @BasePath /
func main() {
	var cfg shared.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	logger, err := shared.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	repository := task.NewRepository(make(map[int]task.Task, 256))

	service := task.NewService(repository)

	controller := task.NewController(service)

	app := BuildRouting("http://"+cfg.Rest.ListenAddress, controller, logger)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}

func BuildRouting(allowOrigins string, controller *task.Controller, logger *zap.SugaredLogger) *fiber.App {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-SomeID",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		AllowOrigins:     allowOrigins,
		MaxAge:           300,
	}))

	apiGroup := app.Group("/v1")
	taskApiGroup := apiGroup.Group("/tasks")

	taskApiGroup.Post("", middlewares.Logging(logger), controller.Create)
	taskApiGroup.Get("", middlewares.Logging(logger), controller.GetAll)
	taskApiGroup.Get("/:id<int>", middlewares.Logging(logger), controller.GetById)
	taskApiGroup.Put("/:id<int>", middlewares.Logging(logger), controller.Update)
	taskApiGroup.Delete("/:id<int>", middlewares.Logging(logger), controller.Delete)

	return app
}
