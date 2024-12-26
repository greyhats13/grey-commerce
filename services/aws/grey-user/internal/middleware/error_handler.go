// Path: grey-user/internal/middleware/error_handler.go

package middleware

import (
	"errors"
	"runtime/debug"
	"strconv"

	"github.com/aws/smithy-go"
	"github.com/gofiber/fiber/v2"
)

// CustomErrorHandler is a centralized error handler for Fiber.
// It ensures correct HTTP status codes for client/server errors,
// logs optional stacktrace on 5xx, and returns a uniform JSON response.
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 for any unrecognized error
	code := fiber.StatusInternalServerError

	// If it's a Fiber error, retrieve the custom status code
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		var ae smithy.APIError
		if errors.As(err, &ae) {
			var errCode = ae.ErrorCode()
			// convert awsStatus to int
			awsStatus, err := strconv.Atoi(errCode)
			if err != nil {
				panic(err)
			}

			debug.PrintStack()
			return c.Status(int(awsStatus)).JSON(fiber.Map{
				"type":    "cloudsdk",
				"message": ae.ErrorMessage(),
			})
		}
	}
	// Example: Log stacktrace if it's a 5xx error
	if code >= 500 {
		// For demonstration, printing stacktrace to the console.
		// You could also log it with your logger if needed.
		debug.PrintStack()
	}

	// Return a JSON response consistent for all errors
	return c.Status(code).JSON(fiber.Map{
		"type":    "app",
		"message": err.Error(),
	})
}
