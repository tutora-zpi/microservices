package buffer

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"time"
)

type BufferedNotification struct {
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	ID        string    `json:"id"`
}

func NewBufferedNotification(dto dto.NotificationDTO) *BufferedNotification {
	return &BufferedNotification{
		Data:      dto.JSON(),
		Timestamp: time.Now(),
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}
