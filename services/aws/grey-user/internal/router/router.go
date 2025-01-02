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
	userGroup.Patch("/:uuid", userHandler.UpdateUser)
	userGroup.Get("/:uuid", userHandler.GetUser)
	userGroup.Delete("/:uuid", userHandler.DeleteUser)
	userGroup.Get("/", userHandler.ListUsers)
}
