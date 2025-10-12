package buffer

import (
	"notification-serivce/internal/domain/dto"
	"time"
)

type BufferedNotification struct {
	// Buffered notifications timestamp
	ID   int64  `json:"id"`
	Data []byte `json:"data"`
}

func NewBufferedNotification(dto dto.NotificationDTO) *BufferedNotification {
	timestamp := time.Now().UTC()

	return &BufferedNotification{
		ID:   timestamp.UnixNano(),
		Data: dto.JSON(),
	}
}

func (b *BufferedNotification) Age(t *time.Time) time.Duration {
	now := time.Now().UTC().UnixNano()

	if t != nil {
		now = t.UTC().UnixNano()
	}

	age := now - b.ID
	return time.Duration(age)
}

func (b *BufferedNotification) AgeNow() time.Duration {
	return b.Age(nil)
}
