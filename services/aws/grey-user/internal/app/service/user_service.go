// Path: grey-user/internal/app/service/user_service.go

package service

import (
	"context"
	errors "grey-user/internal/app"
	"grey-user/internal/app/model"
	"grey-user/internal/app/repository"
	"grey-user/pkg/cache"
	"grey-user/pkg/utils"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, userId string, updateReq map[string]interface{}) (*model.User, error)
	GetUser(ctx context.Context, userId string) (*model.User, error)
	DeleteUser(ctx context.Context, userId string) error
	ListUsers(ctx context.Context, limit int32, lastKey string) ([]model.User, string, error)
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

func (s *userService) UpdateUser(ctx context.Context, userId string, updateReq map[string]interface{}) (*model.User, error) {
	user, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	// Apply updates
	if fn, ok := updateReq["firstname"]; ok {
		user.Firstname = fn.(string)
	}
	if ln, ok := updateReq["lastname"]; ok {
		user.Lastname = ln.(string)
	}
	// Tambahkan field lain jika diperlukan
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	_ = s.cache.Del(ctx, user.UserId)
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, userId string) (*model.User, error) {
	cached, err := s.cache.Get(ctx, userId)
	if err == nil && cached != "" {
		return deserializeUser(cached)
	}
	user, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	serialized, err := serializeUser(user)
	if err != nil {
		// Jika serialisasi gagal, jangan set ke cache
		return user, nil
	}
	_ = s.cache.Set(ctx, userId, serialized, 5*time.Minute)
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userId string) error {
	if err := s.repo.DeleteUser(ctx, userId); err != nil {
		return err
	}
	_ = s.cache.Del(ctx, userId)
	return nil
}

func (s *userService) ListUsers(ctx context.Context, limit int32, lastKey string) ([]model.User, string, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.ListUsers(ctx, limit, lastKey)
}

func serializeUser(u *model.User) (string, error) {
	bytes, err := utils.JSONMarshal(u)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func deserializeUser(s string) (*model.User, error) {
	var user model.User
	err := utils.JSONUnmarshal([]byte(s), &user)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &user, nil
}
