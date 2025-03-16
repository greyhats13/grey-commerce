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

// UserService defines the interface for user-related operations
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

// NewUserService creates a new UserService
func NewUserService(repo repository.UserRepository, cache cache.Cache) UserService {
	return &userService{repo: repo, cache: cache}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	if user.ShopID == "" || user.Email == "" || user.Role == "" ||
		user.Firstname == "" || user.Lastname == "" || user.Gender == "" {
		return errors.ErrInvalidRequest
	}

	return s.repo.CreateUser(ctx, user)
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, userId string, updateReq map[string]interface{}) (*model.User, error) {
	user, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	// Apply updates
	for key, value := range updateReq {
		switch key {
		case "shopId":
			if v, ok := value.(string); ok {
				user.ShopID = v
			}
		case "email":
			if v, ok := value.(string); ok {
				user.Email = v
			}
		case "role":
			if v, ok := value.(string); ok {
				user.Role = v
			}
		case "firstname":
			if v, ok := value.(string); ok {
				user.Firstname = v
			}
		case "lastname":
			if v, ok := value.(string); ok {
				user.Lastname = v
			}
		case "gender":
			if v, ok := value.(string); ok {
				user.Gender = v
			}
		case "birthdate":
			if v, ok := value.(string); ok {
				birthdate, err := model.ParseDate(v)
				if err != nil {
					return nil, errors.ErrInvalidRequest
				}
				user.Birthdate = birthdate
			}
		case "addresses":
			if v, ok := value.([]interface{}); ok {
				addresses, err := model.ParseAddresses(v)
				if err == nil {
					user.Addresses = addresses
				}
			}
		case "phones":
			if v, ok := value.([]interface{}); ok {
				phones, err := model.ParsePhones(v)
				if err == nil {
					user.Phones = phones
				}
			}
		case "paymentMethods":
			if v, ok := value.([]interface{}); ok {
				paymentMethods, err := model.ParsePaymentMethods(v)
				if err == nil {
					user.PaymentMethods = paymentMethods
				}
			}
		case "image":
			if v, ok := value.(map[string]interface{}); ok {
				image, err := model.ParseImage(v)
				if err == nil {
					user.Image = image
				}
			}
		// Add more cases for other fields as needed
		}
	}
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	_ = s.cache.Del(ctx, user.UserId)
	return user, nil
}

// GetUser retrieves a user by userId
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
		return user, nil
	}
	_ = s.cache.Set(ctx, userId, serialized, 5*time.Minute)
	return user, nil
}

// DeleteUser deletes a user by userId
func (s *userService) DeleteUser(ctx context.Context, userId string) error {
	if err := s.repo.DeleteUser(ctx, userId); err != nil {
		return err
	}
	_ = s.cache.Del(ctx, userId)
	return nil
}

// ListUsers lists users with pagination
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
