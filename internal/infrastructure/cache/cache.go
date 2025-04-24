package cache

import (
	"context"
	"time"
)

type BasicCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Del(ctx context.Context, keys ...string) error
}

type ExpirableCache interface {
	BasicCache
	SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error
	Expire(ctx context.Context, key string, ttl time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
}

type MultiCache interface {
	BasicCache
	MGet(ctx context.Context, keys ...string) ([]string, error)
	MSet(ctx context.Context, pairs map[string]string) error
}

type ListCache interface {
	BasicCache
	LPush(ctx context.Context, key string, values ...string) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LPop(ctx context.Context, key string) (string, error)
}

type Counter interface {
	Incr(ctx context.Context, key string) (int64, error)
}

type FullCache interface {
	ExpirableCache
	MultiCache
	ListCache
	Counter
}
