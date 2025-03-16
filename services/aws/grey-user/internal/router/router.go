// Path: internal/router/router.go

package router

import (
	"grey-user/internal/app/handler"
	"grey-user/internal/app/service"
	"grey-user/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, zapLogger logger.Logger, userService service.UserService) {
	userHandler := handler.NewUserHandler(userService)

	v1 := app.Group("/v1")
	userGroup := v1.Group("/user")
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Patch("/:userId", userHandler.UpdateUser)
	userGroup.Get("/:userId", userHandler.GetUser)
	userGroup.Delete("/:userId", userHandler.DeleteUser)
	userGroup.Get("/", userHandler.ListUsers)
}
