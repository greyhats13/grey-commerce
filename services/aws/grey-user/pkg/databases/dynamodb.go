package databases

import (
	"grey-user/internal/config"
    "grey-user/pkg/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBClient represents the DynamoDB client
type DynamoDBClient struct {
	Client *dynamodb.DynamoDB
}

// NewDynamoDBClient creates a DynamoDBClient
func NewDynamoDBClient(cfg *config.Config) (*DynamoDBClient, error) {
    sess, err := session.NewAWSSession(cfg)
    if err != nil {
        return nil, err
    }
    db := dynamodb.New(sess)
    return &DynamoDBClient{Client: db}, nil
}
