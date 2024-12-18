// Path: grey-user/cmd/main.go

package main

import (
	"log"

	"grey-user/internal/config"
	"grey-user/internal/middleware"
	"grey-user/internal/router"
	"grey-user/pkg/database/dynamodb"
	"grey-user/pkg/database/redis"
	"grey-user/pkg/logger"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	dynamoClient, err := dynamodb.NewDynamoDBClient(cfg)
	if err != nil {
		zapLogger.Fatal("failed to init dynamodb client", err)
	}

	redisClient := redis.NewRedisClient(cfg)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.ZapLoggerMiddleware(zapLogger))
	app.Use(middleware.ErrorHandler())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.CompressionMiddleware())
	app.Use(recover.New())

	router.SetupRoutes(app, zapLogger, dynamoClient, redisClient, cfg)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		zapLogger.Fatal("failed to start server", err)
	}
}
