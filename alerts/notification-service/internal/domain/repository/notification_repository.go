package repository

import (
	"context"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/models"
)

type NotificationRepository interface {
	Save(ctx context.Context, n ...models.Notification) ([]*dto.NotificationDTO, error)
	MarkAsDelivered(ctx context.Context, id ...string) error
	Get(ctx context.Context, receiverID string, lastNotificationID *string, limit int) ([]dto.NotificationDTO, error)
	Update(ctx context.Context, fields map[string]any, id string) (*dto.NotificationDTO, error)
	Delete(ctx context.Context, clientID string, ids ...string) error
	Close(ctx context.Context) error
}
