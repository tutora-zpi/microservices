package interfaces

import (
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/requests"
)

type NotificationSerivce interface {
	DeleteNotifications(req *requests.DeleteNotificationsRequest, clientID string) error
	FetchNotifications(req *requests.FetchNotificationsRequest) ([]dto.NotificationDTO, error)
}
