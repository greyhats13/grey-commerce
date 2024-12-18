package main

import (
	"log"

	"github.com/greyhats13/services/aws/grey-user/internal/config"
	"github.com/greyhats13/services/aws/grey-user/internal/middleware"
	"github.com/greyhats13/services/aws/grey-user/internal/router"
	"github.com/greyhats13/services/aws/grey-user/pkg/database/dynamodb"
	"github.com/greyhats13/services/aws/grey-user/pkg/database/redis"
	"github.com/greyhats13/services/aws/grey-user/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init logger
	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	// Init DynamoDB
	dynamoClient, err := dynamodb.NewDynamoDBClient(cfg)
	if err != nil {
		zapLogger.Fatal("failed to init dynamodb client", err)
	}

	// Init Redis
	redisClient := redis.NewRedisClient(cfg)

	app := fiber.New(fiber.Config{
		// Using goccy/go-json for performance
		JSONEncoder: goJSONMarshal,
		JSONDecoder: goJSONUnmarshal,
	})

	// Middlewares
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

func goJSONMarshal(v interface{}) ([]byte, error) {
	return fiber.Marshal(v)
}

func goJSONUnmarshal(data []byte, v interface{}) error {
	return fiber.Unmarshal(data, v)
}
