package types

import "time"

type KV interface {
	Set(key, value string, exp time.Duration)
	Get(key string) *string
}
