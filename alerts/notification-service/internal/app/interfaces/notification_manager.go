package interfaces

import (
	"context"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/flow"
	"time"
)

type NotificationManager interface {
	Subscribe(clientID string) (chan []byte, context.CancelFunc, error)
	Unsubscribe(clientID string)
	Push(dto dto.NotificationDTO) error
	IsClientConnected(clientID string) bool
	FlushBufferedNotification(clientID string, channel chan []byte) int

	EnableBuffering(maxSize int, ttl time.Duration)
	GetBufferedNotifications(clientID string) []*flow.BufferedNotification
}
