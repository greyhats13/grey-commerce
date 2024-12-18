// Path: grey-user/internal/config/config.go

package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AWSRegion     string
	DynamoDBTable string
	RedisAddr     string
	RedisPassword string
	Port          string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config.yaml")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("PORT", "8080")

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file not found
	}

	cfg := &Config{
		AWSRegion:     viper.GetString("AWS_REGION"),
		DynamoDBTable: viper.GetString("DYNAMODB_TABLE"),
		RedisAddr:     viper.GetString("REDIS_ADDR"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
		Port:          viper.GetString("PORT"),
	}

	return cfg, nil
}
