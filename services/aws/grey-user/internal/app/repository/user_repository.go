//Path: grey-user/internal/app/repository/user_repository.go

package repository

import (
	"context"
	"grey-user/internal/app"
	"grey-user/internal/app/model"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, uuid string) (*model.User, error)
	DeleteUser(ctx context.Context, uuid string) error
	ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error)
}

type userRepository struct {
	db    *dynamodb.DynamoDB
	table string
}

func NewUserRepository(db *dynamodb.DynamoDB, table string) UserRepository {
	return &userRepository{
		db:    db,
		table: table,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	user.UUID = uuid.New().String()
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(uuid)"),
	})
	return err
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now().UTC()

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	return err
}

func (r *userRepository) GetUser(ctx context.Context, userID string) (*model.User, error) {
	out, err := r.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(userID)},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, app.ErrNotFound
	}

	var user model.User
	if err := dynamodbattribute.UnmarshalMap(out.Item, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.db.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(userID)},
		},
		ConditionExpression: aws.String("attribute_exists(uuid)"),
	})
	return err
}

func (r *userRepository) ListUsers(ctx context.Context, limit int64, lastKey string) ([]model.User, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.table),
		Limit:     aws.Int64(limit),
	}
	if lastKey != "" {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(lastKey)},
		}
	}

	out, err := r.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, "", err
	}

	users := []model.User{}
	if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &users); err != nil {
		return nil, "", err
	}

	nextKey := ""
	if out.LastEvaluatedKey != nil && out.LastEvaluatedKey["uuid"] != nil {
		nextKey = aws.StringValue(out.LastEvaluatedKey["uuid"].S)
	}

	return users, nextKey, nil
}
