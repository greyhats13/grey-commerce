// Path: grey-user/internal/config/config.go

package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AWSRegion          string
	LocalStackEndpoint string
	AppEnv             string
	DynamoDBTable      string
	RedisAddr          string
	RedisPassword      string
	Port               string
	MySQLUser          string
	MySQLPassword      string
	MySQLHost          string
	MySQLPort          string
	MySQLDBName        string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config.yaml")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("MYSQL_PORT", "3306")

	if err := viper.ReadInConfig(); err != nil {

	}

	cfg := &Config{
		AWSRegion:          viper.GetString("AWS_REGION"),
		LocalStackEndpoint: viper.GetString("LOCALSTACK_ENDPOINT"),
		AppEnv:             viper.GetString("APP_ENV"),
		DynamoDBTable:      viper.GetString("DYNAMODB_TABLE"),
		RedisAddr:          viper.GetString("REDIS_ADDR"),
		RedisPassword:      viper.GetString("REDIS_PASSWORD"),
		Port:               viper.GetString("PORT"),
		MySQLUser:          viper.GetString("MYSQL_USER"),
		MySQLPassword:      viper.GetString("MYSQL_PASSWORD"),
		MySQLHost:          viper.GetString("MYSQL_HOST"),
		MySQLPort:          viper.GetString("MYSQL_PORT"),
		MySQLDBName:        viper.GetString("MYSQL_DB_NAME"),
	}

	return cfg, nil
}
