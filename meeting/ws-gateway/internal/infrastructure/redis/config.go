package redis

import (
	"os"
	"strconv"
	"time"
	"ws-gateway/internal/config"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int

	Timeout time.Duration
}

func NewRedisConfig(timeout time.Duration) *RedisConfig {
	dbStr := os.Getenv(config.REDIS_DB)
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	return &RedisConfig{
		Addr:     os.Getenv(config.REDIS_ADDR),
		Password: os.Getenv(config.REDIS_PASSWORD),
		DB:       db,
		Timeout:  timeout,
	}
}
