package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(addr, password string, db int) (*RedisAdapter, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Redis unavailable, disabling cache: %v\n", err)
		return nil, err
	}
	return &RedisAdapter{client: rdb}, nil
}

// BasicCache
func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
func (r *RedisAdapter) Set(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
func (r *RedisAdapter) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// ExpirableCache
func (r *RedisAdapter) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}
func (r *RedisAdapter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}
func (r *RedisAdapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Counter
func (r *RedisAdapter) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// MultiCache
func (r *RedisAdapter) MGet(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	out := make([]string, len(vals))
	for i, v := range vals {
		if v != nil {
			out[i] = v.(string)
		}
	}
	return out, nil
}

func (r *RedisAdapter) MSet(ctx context.Context, pairs map[string]string) error {
	args := make([]interface{}, 0, len(pairs)*2)
	for k, v := range pairs {
		args = append(args, k, v)
	}
	return r.client.MSet(ctx, args...).Err()
}

// ListCache
func (r *RedisAdapter) LPush(ctx context.Context, key string, values ...string) error {
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}
	return r.client.LPush(ctx, key, args...).Err()
}

func (r *RedisAdapter) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

func (r *RedisAdapter) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}
