package dynamodb

import (
	"grey-user/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDB struct {
	Client *dynamodb.DynamoDB
}

func NewDynamoDBClient(cfg *config.Config) (*DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	db := dynamodb.New(sess)

	// Ensure table exists or create if needed - For real scenario, handle table creation outside
	// Here we assume the table already exists

	return &DynamoDB{Client: db}, nil
}
