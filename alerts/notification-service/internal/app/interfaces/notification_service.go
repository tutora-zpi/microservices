package interfaces

import (
	"notification-serivce/internal/domain/dto"
)

type NotificationSerivce interface {
	DeleteNotifications(req *dto.DeleteNotificationsDTO, clientID string) error
	FetchNotifications(req *dto.FetchNotificationsDTO) ([]dto.NotificationDTO, error)
}
