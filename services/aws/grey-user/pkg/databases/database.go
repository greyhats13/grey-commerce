// File: pkg/databases/database.go
package databases

import (
	"context"
)

type Item map[string]interface{}

type Database interface {
	PutItem(ctx context.Context, item interface{}, conditionExpression *string, tableName string) error
	GetItem(ctx context.Context, key map[string]interface{}, tableName string) (map[string]interface{}, error)
	DeleteItem(ctx context.Context, key map[string]interface{}, tableName string) error
	QueryItems(ctx context.Context, tableName string, limit int64, lastKey string) ([]map[string]interface{}, string, error)
}
