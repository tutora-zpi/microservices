package repository

import (
	"context"
	"signaling-service/internal/domain/enum"
	"signaling-service/internal/domain/models"
	"time"
)

type StatusRepository interface {
	Save(ctx context.Context, userID string, status enum.UserStatus, ttl time.Duration) error
	Get(ctx context.Context, userID string) (*models.Status, error)
	Delete(ctx context.Context, userID string) error
}
