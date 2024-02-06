package core

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisCache struct {
	RDB *redis.Client
}

var ctx = context.Background()

func NewRedisCache(rdb *redis.Client) *RedisCache {
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

func (cache *RedisCache) RemovePrefix(prefix string) error {
	keys, err := cache.RDB.Keys(ctx, prefix+"*").Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err := cache.RDB.Del(ctx, key).Result()
		if err != nil {
			return err
		}
	}
	return nil
}
