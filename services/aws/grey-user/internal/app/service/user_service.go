// Path: internal/app/service/user_service.go
package service

import (
	"context"
	errors "grey-user/internal/app"
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

type userService struct {
	repo  repository.UserRepository
	cache cache.Cache
}

func NewUserService(repo repository.UserRepository, cache cache.Cache) UserService {
	return &userService{repo: repo, cache: cache}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	if user.ShopID == "" || user.Email == "" || user.Role == "" ||
		user.Firstname == "" || user.Lastname == "" || user.Gender == "" {
		return errors.ErrInvalidRequest
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, uuidStr string, updateReq map[string]interface{}) (*model.User, error) {
	user, err := s.repo.GetUser(ctx, uuidStr)
	if err != nil {
		return nil, err
	}
	// apply updates
	if fn, ok := updateReq["firstname"]; ok {
		user.Firstname = fn.(string)
	}
	if ln, ok := updateReq["lastname"]; ok {
		user.Lastname = ln.(string)
	}
	// add more fields if needed
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	_ = s.cache.Del(ctx, user.UUID)
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, uuidStr string) (*model.User, error) {
	cached, err := s.cache.Get(ctx, uuidStr)
	if err == nil && cached != "" {
		return deserializeUser(cached)
	}
	user, err := s.repo.GetUser(ctx, uuidStr)
	if err != nil {
		return nil, err
	}
	serialized := serializeUser(user)
	_ = s.cache.Set(ctx, uuidStr, serialized, 5*time.Minute)
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, uuidStr string) error {
	if err := s.repo.DeleteUser(ctx, uuidStr); err != nil {
		return err
	}
	_ = s.cache.Del(ctx, uuidStr)
	return nil
}

func (s *userService) ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.ListUsers(ctx, limit, lastKey)
}

func serializeUser(u *model.User) string {
	return u.UUID + "|" + u.Email // minimal example; ideally JSON
}

func deserializeUser(s string) (*model.User, error) {
	parts := []rune(s)
	if len(parts) < 2 {
		return nil, errors.ErrNotFound
	}
	// naive
	user := &model.User{
		UUID:  string(parts[0 : len(parts)/2]),
		Email: string(parts[len(parts)/2:]),
	}
	return user, nil
}
