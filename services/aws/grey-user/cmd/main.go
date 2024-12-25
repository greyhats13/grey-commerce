// Path: grey-user/cmd/main.go

package main

import (
    "log"

    "grey-user/internal/app/repository"
    "grey-user/internal/app/service"
    "grey-user/internal/config"
    "grey-user/internal/middleware"
    "grey-user/internal/router"
    "grey-user/pkg/logger"
		"grey-user/pkg/databases"
		"grey-user/pkg/cache"

    "github.com/goccy/go-json"
    "github.com/gofiber/fiber/v2"
)

// This example uses only DynamoDB for the database and Redis for the cache
func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    zapLogger, err := logger.NewZapLogger()
    if err != nil {
        log.Fatalf("failed to init logger: %v", err)
    }

    // Initialize DynamoDB
    dynamoClient, err := databases.NewDynamoDBClient(cfg)
    if err != nil {
        zapLogger.Fatal("failed to create DynamoDB client", err)
    }
    dynamoDatabase := databases.NewDynamoDBDatabase(dynamoClient)
    userRepo := repository.NewDynamoDBUserRepository(dynamoDatabase, cfg.DynamoDBTable)

    // Initialize Redis
    redisClient, err := cache.NewRedisClient(cfg)
    if err != nil {
        zapLogger.Fatal("failed to connect Redis", err)
    }
    redisCache := cache.NewRedisCache(redisClient)

    // Dependency injection for service
    userService := service.NewUserService(userRepo, redisCache)

    app := fiber.New(fiber.Config{
        JSONEncoder: json.Marshal,
        JSONDecoder: json.Unmarshal,
    })

    // Middlewares
    app.Use(middleware.RequestIDMiddleware())
    app.Use(middleware.ZapLoggerMiddleware(zapLogger))
    app.Use(middleware.ErrorHandler())
    app.Use(middleware.CORSMiddleware())
    app.Use(middleware.CompressionMiddleware())
    app.Use(middleware.NewRecoverMiddleware())

    // Routes
    router.SetupRoutes(app, zapLogger, userService)

    port := cfg.Port
    if port == "" {
        port = "8080"
    }
    if err := app.Listen(":" + port); err != nil {
        zapLogger.Fatal("failed to start server", err)
    }
}