package databases

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// dynamoDBDatabase is an implementation of the Database interface for DynamoDB
type dynamoDBDatabase struct {
	client *DynamoDBClient
}

// NewDynamoDBDatabase returns a struct implementing Database interface
func NewDynamoDBDatabase(client *DynamoDBClient) Database {
	return &dynamoDBDatabase{client: client}
}

// PutItem puts or updates an item in a DynamoDB table
func (d *dynamoDBDatabase) PutItem(ctx context.Context, item interface{}, conditionExpression *string, tableName string) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}
	if conditionExpression != nil {
		input.ConditionExpression = conditionExpression
	}
	_, err = d.client.Client.PutItemWithContext(ctx, input)
	return err
}

// GetItem retrieves an item from a DynamoDB table
func (d *dynamoDBDatabase) GetItem(ctx context.Context, key map[string]interface{}, tableName string) (map[string]interface{}, error) {
	keyAV, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return nil, err
	}
	res, err := d.client.Client.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       keyAV,
	})
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, fmt.Errorf("not found")
	}
	out := map[string]interface{}{}
	if err := dynamodbattribute.UnmarshalMap(res.Item, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteItem deletes an item from a DynamoDB table
func (d *dynamoDBDatabase) DeleteItem(ctx context.Context, key map[string]interface{}, tableName string) error {
	keyAV, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return err
	}
	_, err = d.client.Client.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       keyAV,
	})
	return err
}

// QueryItems scans or queries items from a DynamoDB table
func (d *dynamoDBDatabase) QueryItems(ctx context.Context, tableName string, limit int64, lastKey string) ([]map[string]interface{}, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int64(limit),
	}
	if lastKey != "" {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(lastKey)},
		}
	}
	res, err := d.client.Client.ScanWithContext(ctx, input)
	if err != nil {
		return nil, "", err
	}
	items := make([]map[string]interface{}, 0, len(res.Items))
	for _, i := range res.Items {
		tmp := map[string]interface{}{}
		if err := dynamodbattribute.UnmarshalMap(i, &tmp); err == nil {
			items = append(items, tmp)
		}
	}
	next := ""
	if res.LastEvaluatedKey != nil && res.LastEvaluatedKey["uuid"] != nil {
		next = *res.LastEvaluatedKey["uuid"].S
	}
	return items, next, nil
}
