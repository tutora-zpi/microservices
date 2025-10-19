package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, redisConfig RedisConfig, timeout time.Duration) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client := redis.NewClient(opts)

	errCh := make(chan error, 1)
	go func() {
		_, err := client.Ping(ctx).Result()
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		_ = client.Close()
		return nil, fmt.Errorf("redis: connection timeout after %v", timeout)

	case err := <-errCh:
		if err != nil {
			_ = client.Close()
			return nil, fmt.Errorf("redis: failed to ping: %w", err)
		}

		log.Println("Successfully connected to Redis")
		return client, nil
	}
}
