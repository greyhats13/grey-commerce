// Path: grey-user/internal/middleware/logger.go

package middleware

import (
	"errors"
	"grey-user/pkg/logger"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware assigns a request ID to each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := uuid.New().String()
		c.Locals("requestId", id)
		return c.Next()
	}
}

// ZapLoggerMiddleware logs each request using the provided zap-compatible logger
func ZapLoggerMiddleware(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record start time and read the request body size
		start := time.Now()
		requestSize := len(c.Request().Body())

		// Continue to next middleware or actual route
		err := c.Next()
    
		// Retrieve final status code: if err != nil, we check if it's a fiber.Error
		status := c.Response().StatusCode() // Default if no fiber.Error is found

		if err != nil {
			// Attempt to parse the fiber.Error for an accurate status code
			var fe *fiber.Error
			if errors.As(err, &fe) {
				status = fe.Code
			} else {
				// Otherwise, assume it's a server error
				status = fiber.StatusInternalServerError
			}
		}

		// End time for latency
		stop := time.Now()

		// Determine severity
		severity := "INFO"
		switch {
		case status >= 500:
			severity = "ERROR"
		case status >= 400:
			severity = "WARNING"
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

		// If there was an error, override the message with error details
		if err != nil {
			fields = append(fields, logger.Field{Key: "message", Value: err.Error()})

			// Log according to severity
			if severity == "ERROR" {
				log.Error("", fields...)
			} else {
				log.Warn("", fields...)
			}
			return err
		}

		// No error, log as INFO
		log.Info("", fields...)
		return nil
	}
}
