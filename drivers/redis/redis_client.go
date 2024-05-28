package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisClient{
		Client: rdb,
	}
}

func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Del(key string) error {
	return r.Client.Del(ctx, key).Err()
}