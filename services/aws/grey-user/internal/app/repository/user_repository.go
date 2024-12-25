// Path: grey-user/internal/app/repository/user_repository.go

package repository

import (
	"context"
	"grey-user/internal/app/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, uuid string) (*model.User, error)
	DeleteUser(ctx context.Context, uuid string) error
	ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error)
}
