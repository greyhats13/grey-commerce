// Path: internal/app/service/user_service.go

package service

import (
	"context"
	"grey-user/internal/app"
	"grey-user/internal/app/model"
	"grey-user/internal/app/repository"
	"grey-user/pkg/cache"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, uuid string, updateReq map[string]interface{}) (*model.User, error)
	GetUser(ctx context.Context, uuid string) (*model.User, error)
	DeleteUser(ctx context.Context, uuid string) error
	ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error)
}

// userService is the concrete implementation of UserService
type userService struct {
	repo  repository.UserRepository
	cache cache.Cache
}

// NewUserService creates a userService with the required dependencies
func NewUserService(repo repository.UserRepository, c cache.Cache) UserService {
	return &userService{repo: repo, cache: c}
}

// CreateUser delegates user creation to the repository
func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	// Validate essential fields
	if user.ShopID == "" || user.Email == "" || user.Role == "" ||
		user.Firstname == "" || user.Lastname == "" || user.Birthdate.IsZero() ||
		user.Gender == "" || len(user.Addresses) == 0 || len(user.Phones) == 0 {
		return app.ErrInvalidRequest
	}
	return s.repo.CreateUser(ctx, user)
}

// UpdateUser updates user data
func (s *userService) UpdateUser(ctx context.Context, uuidStr string, updateReq map[string]interface{}) (*model.User, error) {
	// Retrieve existing user
	user, err := s.repo.GetUser(ctx, uuidStr)
	if err != nil {
		return nil, err
	}

	// Simple partial updates if found in updateReq
	if fn, ok := updateReq["firstname"]; ok {
		user.Firstname = fn.(string)
	}
	if ln, ok := updateReq["lastname"]; ok {
		user.Lastname = ln.(string)
	}
	if em, ok := updateReq["email"]; ok {
		user.Email = em.(string)
	}
	// More fields can be added

	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.Del(ctx, user.UUID)
	return user, nil
}

// GetUser retrieves user from cache or repository
func (s *userService) GetUser(ctx context.Context, uuidStr string) (*model.User, error) {
	cached, err := s.cache.Get(ctx, uuidStr)
	if err == nil && cached != "" {
		return deserializeUser(cached)
	}
	usr, err := s.repo.GetUser(ctx, uuidStr)
	if err != nil {
		return nil, err
	}

	serialized := serializeUser(usr)
	_ = s.cache.Set(ctx, uuidStr, serialized, 5*time.Minute)
	return usr, nil
}

// DeleteUser removes a user from database and cache
func (s *userService) DeleteUser(ctx context.Context, uuidStr string) error {
	if err := s.repo.DeleteUser(ctx, uuidStr); err != nil {
		return err
	}
	_ = s.cache.Del(ctx, uuidStr)
	return nil
}

// ListUsers delegates to repository
func (s *userService) ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.ListUsers(ctx, limit, lastKey)
}

// Simple serialization to demonstrate caching
func serializeUser(u *model.User) string {
	return u.UUID + "|" + u.Email
}

// Simple deserialization
func deserializeUser(s string) (*model.User, error) {
	r := []rune(s)
	if len(r) < 3 {
		return nil, app.ErrNotFound
	}
	// Arbitrary splitting
	mid := len(r) / 2
	u := &model.User{
		UUID:  string(r[:mid]),
		Email: string(r[mid:]),
	}
	return u, nil
}
