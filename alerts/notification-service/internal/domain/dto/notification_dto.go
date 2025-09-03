package dto

import (
	"encoding/json"
	"notification-serivce/internal/domain/enums"
	"time"
)

type NotificationDTO struct {
	ID              string                 `json:"id"`
	ReceiverID      string                 `json:"receiverID"`
	CreatedAt       time.Time              `json:"createdAt"`
	Type            enums.NotificationType `json:"type"`
	Title           string                 `json:"title"`
	Body            string                 `json:"body"`
	RedirectionLink string                 `json:"redirectionLink"`
	Metadata        map[string]any         `json:"metadata"`
}

func (dto *NotificationDTO) JSON() []byte {
	data, err := json.Marshal(*dto)

	if err != nil {
		return []byte{}
	}

	return data
}
