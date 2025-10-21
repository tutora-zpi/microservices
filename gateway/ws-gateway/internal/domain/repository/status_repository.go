package repository

import (
	"context"
	"time"
	"ws-gateway/internal/domain/enum"
	"ws-gateway/internal/domain/models"
)

type StatusRepository interface {
	Save(ctx context.Context, userID string, status enum.UserStatus, ttl time.Duration) error
	Get(ctx context.Context, userID string) (*models.Status, error)
	Delete(ctx context.Context, userID string) error
}
