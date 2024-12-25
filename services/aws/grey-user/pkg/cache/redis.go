// Path: pkg/cache/redis/redis.go

package cache

import (
	"context"
	"grey-user/internal/config"

	goredis "github.com/redis/go-redis/v9"
)

// RedisClient holds the lower-level redis client
type RedisClient struct {
	Client *goredis.Client
}

// NewRedisClient initializes the RedisClient from config
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{Client: rdb}, nil
}
