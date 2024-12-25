// Path: grey-user/pkg/databases/dynamodb.go

package databases

import (
	"grey-user/internal/config"
	"grey-user/pkg/aws/session"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBClient represents the DynamoDB client for AWS SDK v2
type DynamoDBClient struct {
	Client *dynamodb.Client
}

// NewDynamoDBClient creates a DynamoDBClient using AWS SDK v2
func NewDynamoDBClient(cfg *config.Config) (*DynamoDBClient, error) {
	// Create a new AWS session/config using v2
	awsCfg, err := session.NewAWSSession(cfg)
	if err != nil {
		return nil, err
	}

	// Create a DynamoDB client from the loaded config
	db := dynamodb.NewFromConfig(awsCfg)
	return &DynamoDBClient{Client: db}, nil
}
