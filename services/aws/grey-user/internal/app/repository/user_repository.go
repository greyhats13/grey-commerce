package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/greyhats13/services/aws/grey-user/internal/app"
	"github.com/greyhats13/services/aws/grey-user/internal/app/model"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	item, err := MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
		ConditionExpression: aws.String("attribute_not_exists(uuid)"), 
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now().UTC()

	item, err := MarshalMap(user)
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
	err = UnmarshalMap(out.Item, &user)
	if err != nil {
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
	if err != nil {
		return err
	}
	return nil
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
	err = UnmarshalListOfMaps(out.Items, &users)
	if err != nil {
		return nil, "", err
	}

	nextKey := ""
	if out.LastEvaluatedKey != nil && out.LastEvaluatedKey["uuid"] != nil {
		nextKey = aws.StringValue(out.LastEvaluatedKey["uuid"].S)
	}

	return users, nextKey, nil
}

// MarshalMap and Unmarshal helpers
// For simplicity, use dynamodbattribute utilities
func MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := dynamodbattributeMarshalMap(in)
	if err != nil {
		return nil, err
	}
	return av, nil
}

func UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattributeUnmarshalMap(m, out)
}

func UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattributeUnmarshalListOfMaps(l, out)
}

// We define custom marshal/unmarshal to avoid placeholders
// Using standard AWS SDK methods
func dynamodbattributeMarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute{}.MarshalMap(in)
}
func dynamodbattributeUnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute{}.UnmarshalMap(m, out)
}
func dynamodbattributeUnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute{}.UnmarshalListOfMaps(l, out)
}

// A mini wrapper to handle attribute conversions
type dynamodbattribute struct{}

func (d dynamodbattribute) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	ma, err := dynamodbattributeMarshal(in)
	if err != nil {
		return nil, err
	}
	m, ok := ma.(map[string]*dynamodb.AttributeValue)
	if !ok {
		return nil, fmt.Errorf("unexpected type")
	}
	return m, nil
}

func (d dynamodbattribute) UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattributeUnmarshal(m, out)
}

func (d dynamodbattribute) UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattributeUnmarshal(l, out)
}

// Use dynamodbattribute from aws-sdk-go
// We'll implement these using reflection from the sdk
import "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

func dynamodbattributeMarshal(in interface{}) (interface{}, error) {
	return dynamodbattribute.MarshalMap(in)
}

func dynamodbattributeUnmarshal(av interface{}, out interface{}) error {
	switch v := av.(type) {
	case map[string]*dynamodb.AttributeValue:
		return dynamodbattribute.UnmarshalMap(v, out)
	case []map[string]*dynamodb.AttributeValue:
		return dynamodbattribute.UnmarshalListOfMaps(v, out)
	default:
		return errors.New("unsupported attribute type")
	}
}
