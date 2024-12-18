package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compression"
)

func CompressionMiddleware() fiber.Handler {
	return compression.New()
}
