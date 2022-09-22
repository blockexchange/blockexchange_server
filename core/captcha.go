package core

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCaptchaStore struct {
	rd         *redis.Client
	expiration time.Duration
}

func NewRedisCaptchaStore(rd *redis.Client, expiration time.Duration) *RedisCaptchaStore {
	return &RedisCaptchaStore{rd: rd, expiration: expiration}
}

func (s *RedisCaptchaStore) Set(id string, digits []byte) {
	s.rd.Set(context.Background(), id, digits, s.expiration)
}

func (s *RedisCaptchaStore) Get(id string, clear bool) (digits []byte) {
	res, err := s.rd.Get(context.Background(), id).Result()
	if clear {
		s.Set(id, nil)
	}
	if err == nil {
		return []byte(res)
	}
	return nil
}
