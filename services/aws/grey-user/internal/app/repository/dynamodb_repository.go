// Path: grey-user/internal/app/repository/dynamodb_repository.go

package repository

import (
	"context"
	errors "grey-user/internal/app"
	"grey-user/internal/app/model"
	"grey-user/pkg/databases"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/uuid"
)

// DynamoDBUserRepository is the DynamoDB implementation of UserRepository
type DynamoDBUserRepository struct {
	db    databases.Database
	table string
}

// NewDynamoDBUserRepository returns a new DynamoDBUserRepository
func NewDynamoDBUserRepository(db databases.Database, table string) UserRepository {
	return &DynamoDBUserRepository{
		db:    db,
		table: table,
	}
}

// CreateUser creates a new user in DynamoDB
func (r *DynamoDBUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	user.UUID = uuid.New().String()
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	cond := "attribute_not_exists(uuid)" // ensures we don't overwrite existing user
	return r.db.PutItem(ctx, user, &cond, r.table)
}

// UpdateUser updates an existing user in DynamoDB
func (r *DynamoDBUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	if user.UUID == "" {
		return errors.ErrInvalidRequest
	}
	user.UpdatedAt = time.Now().UTC()
	return r.db.PutItem(ctx, user, nil, r.table)
}

// GetUser retrieves a user by UUID from DynamoDB
func (r *DynamoDBUserRepository) GetUser(ctx context.Context, uuidStr string) (*model.User, error) {
	if uuidStr == "" {
		return nil, errors.ErrInvalidRequest
	}
	key := map[string]interface{}{
		"uuid": uuidStr,
	}
	// db.GetItem returns a map[string]interface{} if the item is found
	res, err := r.db.GetItem(ctx, key, r.table)
	if err != nil {
		// Check if the error is "not found"
		if err.Error() == "not found" {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	// Marshal the res map back to an attribute value
	av, err := attributevalue.MarshalMap(res)
	if err != nil {
		return nil, err
	}

	// Then unmarshal into the user struct
	var user model.User
	if err := attributevalue.UnmarshalMap(av, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser deletes a user from DynamoDB by UUID
func (r *DynamoDBUserRepository) DeleteUser(ctx context.Context, uuidStr string) error {
	if uuidStr == "" {
		return errors.ErrInvalidRequest
	}
	key := map[string]interface{}{
		"uuid": uuidStr,
	}
	return r.db.DeleteItem(ctx, key, r.table)
}

// ListUsers returns a list of users from DynamoDB with pagination
func (r *DynamoDBUserRepository) ListUsers(ctx context.Context, limit int32, lastKey string) ([]model.User, string, error) {
	out, nextKey, err := r.db.QueryItems(ctx, r.table, limit, lastKey)
	if err != nil {
		return nil, "", err
	}

	users := make([]model.User, 0, len(out))
	for _, data := range out {
		av, err := attributevalue.MarshalMap(data)
		if err != nil {
			// skip this record if there's an error
			continue
		}
		var u model.User
		if err := attributevalue.UnmarshalMap(av, &u); err == nil {
			users = append(users, u)
		}
	}

	return users, nextKey, nil
}
