// Path: grey-user/internal/middleware/error_handler.go

package middleware

import (
	"errors"
	"runtime/debug"

	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
	"github.com/gofiber/fiber/v2"
)

func CustomErrorHandler(c *fiber.Ctx, err error) error {
	// Default fallback to 500 for unknown errors
	code := fiber.StatusInternalServerError

	var fe *fiber.Error
	if errors.As(err, &fe) {
		code = fe.Code
	}

	// We'll check if this is an AWS error
	var oe *smithy.OperationError
	var ae smithy.APIError
	var re *http.ResponseError

	// If all three match, it means it's an AWS SDK v2 error
	if errors.As(err, &oe) && errors.As(err, &ae) && errors.As(err, &re) {
		// Convert ae.ErrorCode() to int if it's numeric, or default to 400 if parse fails
		// Often, ae.ErrorCode() might be "400" or "ValidationException".
		// We'll do a small helper parse:
		var awsStatus = 400 // default
		// If ae.ErrorCode() is purely numeric, parse it:
		// If not numeric (like "ValidationException"), we might still keep 400 or handle differently
		// For simplicity, let's keep 400 for all client side errors. You can enhance as needed.

		// Print stacktrace if 4xx or 5xx from AWS
		debug.PrintStack()

		return c.Status(awsStatus).JSON(fiber.Map{
			"type":    "AWS",
			"message": ae.ErrorMessage(),
		})
	}

	// If it's not an AWS error, we consider it an app error
	if code >= 500 {
		// Print stacktrace only for internal server errors
		debug.PrintStack()
	}
	return c.Status(code).JSON(fiber.Map{
		"type":    "app",
		"message": err.Error(),
	})
}