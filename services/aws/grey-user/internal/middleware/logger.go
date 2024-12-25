// Path: grey-user/internal/middleware/logger.go

package middleware

import (
	"grey-user/pkg/logger"
	"strconv"
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
		// Record start time and read the request body size
		start := time.Now()
		requestSize := len(c.Request().Body())

		// Continue to next middleware or actual route
		err := c.Next()

		// End time for latency
		stop := time.Now()

		// Retrieve the final status code set by handlers
		status := c.Response().StatusCode()

		// Figure out severity
		severity := "INFO"
		if status >= 400 && status < 500 {
			severity = "WARNING"
		} else if status >= 500 {
			severity = "ERROR"
		}

		// requestId from locals
		requestId := c.Locals("requestId")
		if requestId == nil {
			requestId = ""
		}

		// responseSize
		responseSize := len(c.Response().Body())

		httpRequest := map[string]interface{}{
			"requestMethod": c.Method(),
			"requestUrl":    c.OriginalURL(),
			"requestSize":   strconv.Itoa(requestSize),
			"status":        status,
			"responseSize":  strconv.Itoa(responseSize),
			"userAgent":     c.Get("User-Agent"),
			"remoteIp":      c.IP(),
			"referer":       c.Get("Referer"),
			"latency":       stop.Sub(start).String(),
			"protocol":      c.Protocol(),
		}

		fields := []logger.Field{
			{Key: "timestamp", Value: time.Now().Format(time.RFC3339)},
			{Key: "severity", Value: severity},
			{Key: "message", Value: "request completed"},
			{Key: "requestId", Value: requestId},
			{Key: "httpRequest", Value: httpRequest},
		}

		if err != nil {
			// Add error details to log fields
			fields = append(fields,
				logger.Field{Key: "message", Value: err.Error()},
			)
			if severity == "ERROR" {
				log.Error("", fields...)
			} else if severity == "WARNING" {
				log.Warn("", fields...)
			} else {
				log.Info("", fields...)
			}
			return err
		}

		log.Info("", fields...)
		return nil
	}
}
