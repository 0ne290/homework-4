package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "homework-4/docs"
	"homework-4/internal"
	"homework-4/internal/middlewares"
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
	var cfg internal.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	logger, err := NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := NewPostgresConnectionPool(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(err)
	}

	unitOfWork := task.NewPostgresUnitOfWork(pool)

	timeProvider := internal.NewRealTimeProvider()
	uuidProvider := internal.NewRealUuidProvider()

	service := task.NewRealService(unitOfWork, timeProvider, uuidProvider)

	controller := task.NewController(service)

	app := BuildRouting("http://"+cfg.Rest.ListenAddress, controller, logger)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err = app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(err)
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}

func NewLogger(level string) (*zap.SugaredLogger, error) {
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	logger, err := zap.Config{
		Level:       logLevel,
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			TimeKey:    "timestamp",
			EncodeTime: zapcore.RFC3339NanoTimeEncoder,
		},
		DisableStacktrace: true,
	}.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func NewPostgresConnectionPool(ctx context.Context, cfg internal.PostgreSQL) (*pgxpool.Pool, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
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
	taskApiGroup.Get("/:uuid<guid>", middlewares.Logging(logger), controller.GetByUuid)
	taskApiGroup.Put("/:uuid<guid>", middlewares.Logging(logger), controller.Update)
	taskApiGroup.Delete("/:uuid<guid>", middlewares.Logging(logger), controller.Delete)

	return app
}
