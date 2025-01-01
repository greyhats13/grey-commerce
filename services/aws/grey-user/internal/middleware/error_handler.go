// Path: grey-user/internal/middleware/error_handler.go

package middleware

import (
	"errors"
	"runtime/debug"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
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
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			// Return a JSON response consistent for all errors
			debug.PrintStack()
			return c.Status(re.HTTPStatusCode()).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
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
		"success": false,
		"message": err.Error(),
	})
}
