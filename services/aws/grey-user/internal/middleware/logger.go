// Path: services/aws/grey-user/internal/middleware/logger.go

package middleware

import (
	"services/aws/grey-user/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := uuid.New().String()
		c.Locals("requestId", id)
		return c.Next()
	}
}

func ZapLoggerMiddleware(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		requestId := c.Locals("requestId")
		if requestId == nil {
			requestId = ""
		}

		httpRequest := map[string]interface{}{
			"requestMethod": c.Method(),
			"requestUrl":    c.OriginalURL(),
			"requestSize":   "-1",
			"status":        c.Response().StatusCode(),
			"responseSize":  "-1",
			"userAgent":     c.Get("User-Agent"),
			"remoteIp":      c.IP(),
			"referer":       c.Get("Referer"),
			"latency":       stop.Sub(start).String(),
			"protocol":      c.Protocol(),
		}

		fields := []logger.Field{
			{Key: "timestamp", Value: time.Now().Format(time.RFC3339)},
			{Key: "severity", Value: "INFO"},
			{Key: "message", Value: "request completed"},
			{Key: "requestId", Value: requestId},
			{Key: "httpRequest", Value: httpRequest},
		}

		if err != nil {
			fields = append(fields,
				logger.Field{Key: "severity", Value: "ERROR"},
				logger.Field{Key: "message", Value: err.Error()})
			log.Error("", fields...)
			return err
		}

		log.Info("", fields...)
		return nil
	}
}
