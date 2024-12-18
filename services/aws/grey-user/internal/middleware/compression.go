// Path: grey-user/internal/middleware/compression.go

package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func CompressionMiddleware() fiber.Handler {
	return compress.New()
}
