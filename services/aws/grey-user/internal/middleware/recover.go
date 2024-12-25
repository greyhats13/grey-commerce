// Path: grey-user/internal/middleware/recover.go

package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RecoverMiddleware() fiber.Handler {
	return recover.New()
}
