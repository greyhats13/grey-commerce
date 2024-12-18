package service

import (
	"context"
	"time"

	"github.com/greyhats13/services/aws/grey-user/internal/app"
	"github.com/greyhats13/services/aws/grey-user/internal/app/model"
	"github.com/greyhats13/services/aws/grey-user/internal/app/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, uuid string, updateReq map[string]interface{}) (*model.User, error)
	GetUser(ctx context.Context, uuid string) (*model.User, error)
	DeleteUser(ctx context.Context, uuid string) error
	ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	// Basic validation
	if user.ShopID == "" || user.Email == "" || user.Role == "" || user.Firstname == "" || user.Lastname == "" || user.Birthdate.IsZero() || user.Gender == "" || len(user.Addresses) == 0 || len(user.Phones) == 0 {
		return app.ErrInvalidRequest
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, uuid string, updateReq map[string]interface{}) (*model.User, error) {
	user, err := s.repo.GetUser(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if val, ok := updateReq["shopId"]; ok {
		user.ShopID = val.(string)
	}
	if val, ok := updateReq["email"]; ok {
		user.Email = val.(string)
	}
	if val, ok := updateReq["role"]; ok {
		user.Role = val.(string)
	}
	if val, ok := updateReq["firstname"]; ok {
		user.Firstname = val.(string)
	}
	if val, ok := updateReq["lastname"]; ok {
		user.Lastname = val.(string)
	}
	if val, ok := updateReq["birthdate"]; ok {
		// birthdate should be a string in "2006-01-02" format or RFC3339
		strVal := val.(string)
		bd, err := time.Parse("2006-01-02", strVal)
		if err != nil {
			return nil, app.ErrInvalidRequest
		}
		user.Birthdate = bd
	}
	if val, ok := updateReq["gender"]; ok {
		user.Gender = model.Gender(val.(string))
	}
	if val, ok := updateReq["addresses"]; ok {
		user.Addresses = toAddresses(val)
	}
	if val, ok := updateReq["phones"]; ok {
		user.Phones = toPhones(val)
	}
	if val, ok := updateReq["image"]; ok {
		user.Image = toImage(val)
	}
	if val, ok := updateReq["paymentMethods"]; ok {
		user.PaymentMethods = toPaymentMethods(val)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, uuid string) (*model.User, error) {
	return s.repo.GetUser(ctx, uuid)
}

func (s *userService) DeleteUser(ctx context.Context, uuid string) error {
	return s.repo.DeleteUser(ctx, uuid)
}

func (s *userService) ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.ListUsers(ctx, limit, lastKey)
}

// helper casting functions
func toAddresses(val interface{}) []model.Address {
	arr, ok := val.([]interface{})
	if !ok {
		return nil
	}
	addrs := []model.Address{}
	for _, v := range arr {
		m := v.(map[string]interface{})
		addrs = append(addrs, model.Address{
			Type:        toString(m["type"]),
			Address:     toString(m["address"]),
			Subdistrict: toString(m["subdistrict"]),
			District:    toString(m["district"]),
			City:        toString(m["city"]),
			Province:    toString(m["province"]),
			Country:     toString(m["country"]),
			PostalCode:  toString(m["postalCode"]),
		})
	}
	return addrs
}

func toPhones(val interface{}) []model.Phone {
	arr, ok := val.([]interface{})
	if !ok {
		return nil
	}
	phones := []model.Phone{}
	for _, v := range arr {
		m := v.(map[string]interface{})
		phones = append(phones, model.Phone{
			Type:   toString(m["type"]),
			Number: toString(m["number"]),
		})
	}
	return phones
}

func toImage(val interface{}) *model.Image {
	m, ok := val.(map[string]interface{})
	if !ok {
		return nil
	}
	return &model.Image{
		Name: toString(m["name"]),
		URL:  toString(m["url"]),
	}
}

func toPaymentMethods(val interface{}) []model.PaymentMethod {
	arr, ok := val.([]interface{})
	if !ok {
		return nil
	}
	pms := []model.PaymentMethod{}
	for _, v := range arr {
		m := v.(map[string]interface{})
		pms = append(pms, model.PaymentMethod{
			Type:   toString(m["type"]),
			Name:   toString(m["name"]),
			Number: toString(m["number"]),
		})
	}
	return pms
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return v.(string)
}
