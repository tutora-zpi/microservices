package interfaces

import "notification-serivce/internal/domain/dto"

type NotificationPublisher interface {
	Push(dto dto.NotificationDTO) error
	Subscribe(clientID string) (chan []byte, error)
	Unsubscribe(clientID string) error
}
