// Path: services/aws/grey-user/internal/middleware/error_handler.go
package middleware

import "github.com/gofiber/fiber/v2"

func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
