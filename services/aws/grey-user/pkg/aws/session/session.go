// Path: grey-user/pkg/aws/session/session.go

package session

import (
    "grey-user/internal/config"

    "github.com/aws/aws-sdk-go/aws"
    awssession "github.com/aws/aws-sdk-go/aws/session"
)

// NewAWSSession creates and returns an AWS session
func NewAWSSession(cfg *config.Config) (*awssession.Session, error) {
    awsConfig := &aws.Config{
        Region: aws.String(cfg.AWSRegion),
    }

    // Use LocalStack endpoint if AppEnv is "local"
    if cfg.AppEnv == "local" {
        awsConfig.Endpoint = aws.String(cfg.LocalStackEndpoint)
    }

    sess, err := awssession.NewSession(awsConfig)
    if err != nil {
        return nil, err
    }
    return sess, nil
}