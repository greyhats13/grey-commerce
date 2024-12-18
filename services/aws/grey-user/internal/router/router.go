// Path: services/aws/grey-user/internal/router/router.go

package router

import (
	"services/aws/grey-user/internal/app/handler"
	"services/aws/grey-user/internal/app/repository"
	"services/aws/grey-user/internal/app/service"
	"services/aws/grey-user/internal/config"
	"services/aws/grey-user/pkg/database/dynamodb"
	"services/aws/grey-user/pkg/logger"

	// Import the redis client from github.com/redis/go-redis/v9
	"github.com/redis/go-redis/v9"

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
