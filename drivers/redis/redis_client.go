package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: "default",
		DB:       0,
	})

	return &RedisClient{
		Client: rdb,
	}
}

func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		fmt.Printf("Redis SET error: %v\n", err)
	}
	return err
}

func (r *RedisClient) Get(key string) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("Redis GET error: %v\n", err)
	}
	return result, err
}

func (r *RedisClient) Del(key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		fmt.Printf("Redis DEL error: %v\n", err)
	}
	return err
}
