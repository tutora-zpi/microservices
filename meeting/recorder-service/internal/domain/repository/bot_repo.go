package repository

import (
	"context"
	"time"
)

type BotRepository interface {
	TryAdd(ctx context.Context, roomID, botID string, ttl time.Duration) error
	Delete(ctx context.Context, roomID string) error
}
