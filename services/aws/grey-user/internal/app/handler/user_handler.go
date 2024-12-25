// Path: grey-user/internal/app/handler/user_handler.go

package handler

import (
	"net/http"

	errors "grey-user/internal/app"
	"grey-user/internal/app/model"
	"grey-user/internal/app/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// validate is a singleton validator instance
var validate = validator.New()

type UserHandler struct {
	service service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser handles creation of a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		// Return 400 if body parsing fails
		return fiber.NewError(http.StatusBadRequest, errors.ErrFailedToParse.Error())
	}

	// Validate struct fields
	if err := validate.Struct(user); err != nil {
		// You can either return a generic message
		// or parse the validation error for detailed info
		return fiber.NewError(http.StatusBadRequest, errors.ErrFailedToValidate.Error())
	}

	if err := h.service.CreateUser(c.Context(), &user); err != nil {
		if err == errors.ErrInvalidRequest {
			// Possibly something missing or invalid in the user
			return fiber.NewError(http.StatusBadRequest, errors.ErrFailedToParse.Error())
		} else if err == errors.ErrFailedToValidate {
			return fiber.NewError(http.StatusBadRequest, errors.ErrFailedToValidate.Error())
		}
		// For any other errors, treat it as internal
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// On success, return HTTP 201
	return c.Status(http.StatusCreated).JSON(user)
}

// UpdateUser handles updating of an existing user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}

	var updateReq map[string]interface{}
	if err := c.BodyParser(&updateReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, errors.ErrInvalidRequest.Error())
	}

	user, err := h.service.UpdateUser(c.Context(), uuidParam, updateReq)
	if err != nil {
		if err == errors.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(user)
}

// GetUser retrieves a user by UUID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}

	user, err := h.service.GetUser(c.Context(), uuidParam)
	if err != nil {
		if err == errors.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(user)
}

// DeleteUser deletes a user by UUID
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}

	if err := h.service.DeleteUser(c.Context(), uuidParam); err != nil {
		if err == errors.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

// ListUsers is just an example; feel free to implement
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	return fiber.NewError(http.StatusNotImplemented, "list users not implemented")
}
