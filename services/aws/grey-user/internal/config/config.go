package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	DynamoDBTable      string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	Port               string
}

func LoadConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		// not fatal if .env is missing, env vars might be set on system
	}

	redisDB := 0

	return &Config{
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DynamoDBTable:      os.Getenv("DYNAMODB_TABLE"),
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDB:            redisDB,
		Port:               os.Getenv("PORT"),
	}, nil
}
