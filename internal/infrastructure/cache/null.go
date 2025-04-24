package cache

import (
	"context"
	"time"
)

type NullCache struct{}

func NewNullAdapter() *NullCache {
	return &NullCache{}
}

// BasicCache
func (n *NullCache) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}
func (n *NullCache) Set(ctx context.Context, key, value string) error {
	return nil
}
func (n *NullCache) Del(ctx context.Context, keys ...string) error {
	return nil
}

// ExpirableCache
func (n *NullCache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	return nil
}
func (n *NullCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}
func (n *NullCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	// TTL -1 signals key does not exist or no TTL
	return -1, nil
}

// MultiCache
func (n *NullCache) MGet(ctx context.Context, keys ...string) ([]string, error) {
	// return slice of empty strings matching length
	res := make([]string, len(keys))
	for i := range res {
		res[i] = ""
	}
	return res, nil
}
func (n *NullCache) MSet(ctx context.Context, pairs map[string]string) error {
	return nil
}

// Counter
func (r *NullCache) Incr(ctx context.Context, key string) (int64, error) {
	return 1, nil
}

// ListCache
func (n *NullCache) LPush(ctx context.Context, key string, values ...string) error {
	return nil
}
func (n *NullCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return []string{}, nil
}
func (n *NullCache) LPop(ctx context.Context, key string) (string, error) {
	return "", nil
}
