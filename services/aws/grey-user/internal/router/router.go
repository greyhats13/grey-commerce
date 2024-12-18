package router

import (
	"github.com/greyhats13/services/aws/grey-user/internal/app/handler"
	"github.com/greyhats13/services/aws/grey-user/internal/app/repository"
	"github.com/greyhats13/services/aws/grey-user/internal/app/service"
	"github.com/greyhats13/services/aws/grey-user/internal/config"

	"github.com/greyhats13/services/aws/grey-user/pkg/database/dynamodb"
	"github.com/greyhats13/services/aws/grey-user/pkg/database/redis"
	"github.com/greyhats13/services/aws/grey-user/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, zapLogger logger.Logger, dynamoClient *dynamodb.DynamoDB, redisClient *redis.Client, cfg *config.Config) {
	userRepo := repository.NewUserRepository(dynamoClient.Client, cfg.DynamoDBTable)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, redisClient)

	v1 := app.Group("/v1")

	userGroup := v1.Group("/user")
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Patch("/:uuid", userHandler.UpdateUser)
	userGroup.Get("/:uuid", userHandler.GetUser)
	userGroup.Delete("/:uuid", userHandler.DeleteUser)
	userGroup.Get("/", userHandler.ListUsers)
}
