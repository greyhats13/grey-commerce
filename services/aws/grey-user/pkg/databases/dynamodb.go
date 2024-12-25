package databases

import (
	"grey-user/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBClient represents the DynamoDB client
type DynamoDBClient struct {
	Client *dynamodb.DynamoDB
}

// NewDynamoDBClient creates a DynamoDBClient
func NewDynamoDBClient(cfg *config.Config) (*DynamoDBClient, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
	})
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)
	return &DynamoDBClient{Client: db}, nil
}
