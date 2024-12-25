// Path: pkg/cache/redis/redis.go

package cache

import (
	"context"
	"time"
)

// RedisCache is a struct that implements the Cache interface using Redis
type RedisCache struct {
	client *RedisClient
}

// NewRedisCache returns a Cache implementation that uses Redis
func NewRedisCache(client *RedisClient) Cache {
	return &RedisCache{client: client}
}

// Get retrieves a value by key from redis
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Client.Get(ctx, key).Result()
}

// Set saves a value with expiration in redis
func (r *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Client.Set(ctx, key, value, expiration).Err()
}

// Del deletes a key from redis
func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Client.Del(ctx, key).Err()
}
