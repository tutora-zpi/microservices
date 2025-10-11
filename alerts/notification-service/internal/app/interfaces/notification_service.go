package interfaces

import (
	"context"
	"notification-serivce/internal/domain/dto"
)

type NotificationSerivce interface {
	DeleteNotifications(ctx context.Context, req *dto.DeleteNotificationsDTO, clientID string) error
	FetchNotifications(ctx context.Context, req *dto.FetchNotificationsDTO) ([]dto.NotificationDTO, error)
}
