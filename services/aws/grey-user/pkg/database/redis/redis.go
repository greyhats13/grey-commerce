// Path: services/aws/grey-user/pkg/database/redis/redis.go

package redis

import (
	"context"
	"services/aws/grey-user/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	// Test connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}
	return rdb
}
