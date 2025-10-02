package dto

import (
	"encoding/json"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/metadata"
)

type NotificationDTO struct {
	ID              string                 `json:"id"`
	Receiver        UserDTO                `json:"receiver"`
	Sender          UserDTO                `json:"sender"`
	CreatedAt       int64                  `json:"createdAt"`
	Type            enums.NotificationType `json:"type"`
	Title           string                 `json:"title"`
	Body            string                 `json:"body"`
	RedirectionLink string                 `json:"redirectionLink"`
	Metadata        map[metadata.Key]any   `json:"metadata"`
}

func (dto *NotificationDTO) JSON() []byte {
	data, err := json.Marshal(*dto)

	if err != nil {
		return []byte{}
	}

	return data
}

func (dto *NotificationDTO) AppendTitle(title string) *NotificationDTO {
	dto.Title = title

	return dto
}

func (dto *NotificationDTO) AppendBody(body string) *NotificationDTO {
	dto.Body = body

	return dto
}
