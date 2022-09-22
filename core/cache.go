package core

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(key string, value []byte, expire time.Duration) error
	Get(key string) ([]byte, error)
	Remove(key string) error
}

// No-op

type NoOpCache struct{}

func (cache *NoOpCache) Set(key string, value []byte, expire time.Duration) error {
	return nil
}

func (cache *NoOpCache) Get(key string) ([]byte, error) {
	return nil, nil
}

func (cache *NoOpCache) Remove(key string) error {
	return nil
}

func NewNoOpCache() Cache {
	return &NoOpCache{}
}

// Redis

type RedisCache struct {
	RDB *redis.Client
}

var ctx = context.Background()

func NewRedisCache(rdb *redis.Client) Cache {
	return &RedisCache{RDB: rdb}
}

func (cache *RedisCache) Set(key string, value []byte, expire time.Duration) error {
	return cache.RDB.Set(ctx, key, string(value), expire).Err()
}

func (cache *RedisCache) Get(key string) ([]byte, error) {
	cmd := cache.RDB.Get(ctx, key)

	data, err := cmd.Result()
	if err == redis.Nil {
		return nil, nil
	}

	return []byte(data), nil
}

func (cache *RedisCache) Remove(key string) error {
	cmd := cache.RDB.Del(ctx, key)
	_, err := cmd.Result()
	return err
}
