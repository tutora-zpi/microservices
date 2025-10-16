package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, redisConfig RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	}

	client := redis.NewClient(opts)

	pong, err := client.Ping(ctx).Result()
	if err != nil || pong != "PONG" {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	log.Println("Successfully connected to Redis:", pong)

	return client, nil
}
