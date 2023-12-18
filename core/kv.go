package core

import (
	"blockexchange/types"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

// redis

type RedisKV struct {
	rc *redis.Client
}

func NewRedisKV(rc *redis.Client) types.KV {
	return &RedisKV{rc: rc}
}

func (l *RedisKV) Set(key, value string, exp time.Duration) {
	l.rc.Set(context.Background(), key, value, exp)
}

func (l *RedisKV) Get(key string) *string {
	s, err := l.rc.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}
	return &s
}

// local

type LocalKV struct {
	c *expirable.LRU[string, string]
}

func (l *LocalKV) Set(key, value string, exp time.Duration) {
	l.c.Add(key, value)
}

func (l *LocalKV) Get(key string) *string {
	v, ok := l.c.Get(key)
	if !ok {
		return nil
	}
	return &v
}

func NewLocalKV() types.KV {
	return &LocalKV{
		c: expirable.NewLRU[string, string](256, nil, time.Minute),
	}
}
