// Path: grey-user/internal/app/handler/user_handler.go

package handler

import (
	"grey-user/internal/app"
	"grey-user/internal/app/model"
	"grey-user/internal/app/service"
	"net/http"

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
		return fiber.NewError(http.StatusBadRequest, app.ErrFailedToParse.Error())
	}

	// Validate struct fields
	if err := validate.Struct(user); err != nil {
		return fiber.NewError(http.StatusBadRequest, app.ErrFailedToValidate.Error())
	}

	if err := h.service.CreateUser(c.Context(), &user); err != nil {
		if err == app.ErrInvalidRequest {
			return fiber.NewError(http.StatusBadRequest, app.ErrFailedToParse.Error())
		} else if err == app.ErrFailedToValidate {
			return fiber.NewError(http.StatusBadRequest, app.ErrFailedToValidate.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusCreated).JSON(user)
}

// UpdateUser handles updating of an existing user
// We changed to match the service signature: UpdateUser(ctx context.Context, uuid string, updateReq map[string]interface{})
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}

	// We will parse into a map instead of a full User struct
	var updateReq map[string]interface{}
	if err := c.BodyParser(&updateReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, app.ErrInvalidRequest.Error())
	}

	user, err := h.service.UpdateUser(c.Context(), uuidParam, updateReq)
	if err != nil {
		if err == app.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(user)
}

// GetUser retrieves a user by UUID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}
	user, err := h.service.GetUser(c.Context(), uuid)
	if err != nil {
		if err == app.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(user)
}

// DeleteUser deletes a user by UUID
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}
	if err := h.service.DeleteUser(c.Context(), uuid); err != nil {
		if err == app.ErrNotFound {
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
