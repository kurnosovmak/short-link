package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
}

type KVContract interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error)
}

func New(addr, password string, db int) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	client := Redis{client: rdb}

	return &client
}

func (redis *Redis) Get(ctx context.Context, key string) (string, error) {
	return redis.client.Get(ctx, key).Result()
}

func (redis *Redis) Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error) {
	return redis.client.Set(ctx, key, value, expiration).Result()
}
