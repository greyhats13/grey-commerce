// Path: grey-user/pkg/aws/session/session.go

package session

import (
	"context"
	"errors"
	"grey-user/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewAWSSession returns an aws.Config object configured for AWS SDK v2
// If LocalStackEndpoint is provided, it will use it as a custom endpoint
func NewAWSSession(cfg *config.Config) (aws.Config, error) {
	var customResolver aws.EndpointResolverWithOptionsFunc

	// If localstack endpoint is set, override the endpoint for DynamoDB
	if cfg.LocalStackEndpoint != "" {
		customResolver = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           cfg.LocalStackEndpoint,
						SigningRegion: cfg.AWSRegion,
					}, nil
				}
				return aws.Endpoint{}, errors.New("unknown endpoint requested")
			},
		)
	}

	loadOpts := []func(*sdkConfig.LoadOptions) error{
		sdkConfig.WithRegion(cfg.AWSRegion),
	}

	// If it has a custom resolver, add it to the load options
	if customResolver != nil {
		loadOpts = append(loadOpts, sdkConfig.WithEndpointResolverWithOptions(customResolver))
	}

	// Load the shared config
	awsCfg, err := sdkConfig.LoadDefaultConfig(context.TODO(), loadOpts...)
	if err != nil {
		return aws.Config{}, err
	}

	return awsCfg, nil
}
