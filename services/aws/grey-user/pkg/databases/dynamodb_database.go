// Path: grey-user/pkg/databases/dynamodb_database.go

package databases

import (
	"context"
	"fmt"

	client "grey-user/pkg/aws/dynamodb"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// dynamoDBDatabase is an implementation of the Database interface for AWS SDK v2
type dynamoDBDatabase struct {
	client *client.DynamoDBClient
}

// NewDynamoDBDatabase returns a struct implementing the Database interface
func NewDynamoDBDatabase(client *client.DynamoDBClient) Database {
	return &dynamoDBDatabase{client: client}
}

// PutItem puts or updates an item in a DynamoDB table
func (d *dynamoDBDatabase) PutItem(ctx context.Context, item interface{}, conditionExpression *string, tableName string) error {
	// Marshal the Go struct (or map) into a dynamodb attribute value map
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: awsString(tableName),
		Item:      av,
	}

	if conditionExpression != nil {
		input.ConditionExpression = conditionExpression
	}

	// AWS SDK v2 uses PutItem(ctx, input) instead of PutItemWithContext
	_, err = d.client.Client.PutItem(ctx, input)
	return err
}

// GetItem retrieves an item from a DynamoDB table
func (d *dynamoDBDatabase) GetItem(ctx context.Context, key map[string]interface{}, tableName string) (map[string]interface{}, error) {
	// Marshal the 'key' map into a dynamodb attribute value map
	keyAV, err := attributevalue.MarshalMap(key)
	if err != nil {
		return nil, err
	}

	// Perform the GetItem request
	res, err := d.client.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: awsString(tableName),
		Key:       keyAV,
	})
	if err != nil {
		return nil, err
	}

	if res.Item == nil {
		return nil, fmt.Errorf("not found")
	}

	// Unmarshal it into a generic map
	out := map[string]interface{}{}
	if err := attributevalue.UnmarshalMap(res.Item, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteItem deletes an item from a DynamoDB table
func (d *dynamoDBDatabase) DeleteItem(ctx context.Context, key map[string]interface{}, tableName string) error {
	keyAV, err := attributevalue.MarshalMap(key)
	if err != nil {
		return err
	}

	_, err = d.client.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: awsString(tableName),
		Key:       keyAV,
	})

	return err
}

// QueryItems scans items from a DynamoDB table (renamed from Query to Scan for brevity)
func (d *dynamoDBDatabase) QueryItems(ctx context.Context, tableName string, limit int32, lastKey string) ([]map[string]interface{}, string, error) {
	// Using scan instead of query for simplicity
	input := &dynamodb.ScanInput{
		TableName: awsString(tableName),
		Limit:     &limit,
	}

	// If lastKey is provided, set the ExclusiveStartKey
	if lastKey != "" {
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: lastKey},
		}
	}

	// Perform the Scan
	res, err := d.client.Client.Scan(ctx, input)
	if err != nil {
		return nil, "", err
	}

	// Convert each item from AttributeValue map to Go map
	items := make([]map[string]interface{}, 0, len(res.Items))
	for _, i := range res.Items {
		tmp := map[string]interface{}{}
		if err := attributevalue.UnmarshalMap(i, &tmp); err == nil {
			items = append(items, tmp)
		}
	}

	// Get the next key if present
	next := ""
	if len(res.LastEvaluatedKey) > 0 {
		if val, ok := res.LastEvaluatedKey["userId"].(*types.AttributeValueMemberS); ok {
			next = val.Value
		}
	}

	return items, next, nil
}

func awsString(val string) *string {
	return &val
}
