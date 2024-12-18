// Path: services/aws/grey-user/internal/app/handler/user_handler.go

package handler

import (
	"net/http"
	"services/aws/grey-user/internal/app"
	"services/aws/grey-user/internal/app/model"
	"services/aws/grey-user/internal/app/service"
	"services/aws/grey-user/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	service     service.UserService
	redisClient *redis.Client
}

func NewUserHandler(service service.UserService, redisClient *redis.Client) *UserHandler {
	return &UserHandler{
		service:     service,
		redisClient: redisClient,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	err := h.service.CreateUser(c.Context(), &user)
	if err != nil {
		if err == app.ErrInvalidRequest {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, "failed to create user")
	}

	return c.Status(http.StatusCreated).JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}

	var updateReq map[string]interface{}
	if err := c.BodyParser(&updateReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	user, err := h.service.UpdateUser(c.Context(), uuid, updateReq)
	if err != nil {
		if err == app.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		} else if err == app.ErrInvalidRequest {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, "failed to update user")
	}

	h.redisClient.Del(c.Context(), user.UUID)
	return c.Status(http.StatusOK).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}

	// Check cache
	cacheKey := uuid
	cachedVal, err := h.redisClient.Get(c.Context(), cacheKey).Result()
	if err == nil && cachedVal != "" {
		var cachedUser model.User
		if err := utils.JSONUnmarshal([]byte(cachedVal), &cachedUser); err == nil {
			return c.JSON(cachedUser)
		}
	}

	user, err := h.service.GetUser(c.Context(), uuid)
	if err != nil {
		if err == app.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, "failed to get user")
	}

	b, _ := utils.JSONMarshal(user)
	h.redisClient.Set(c.Context(), cacheKey, string(b), 5*time.Minute)

	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return fiber.NewError(http.StatusBadRequest, "uuid is required")
	}
	err := h.service.DeleteUser(c.Context(), uuid)
	if err != nil {
		if err == app.ErrNotFound {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, "failed to delete user")
	}

	h.redisClient.Del(c.Context(), uuid)
	return c.SendStatus(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	lastKey := c.Query("lastKey", "")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}

	users, nextKey, err := h.service.ListUsers(c.Context(), limit, lastKey)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to list users")
	}

	resp := map[string]interface{}{
		"users":   users,
		"nextKey": nextKey,
	}
	return c.JSON(resp)
}
