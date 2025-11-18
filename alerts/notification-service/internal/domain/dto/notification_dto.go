package dto

import (
	"encoding/json"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/metadata"
)

type NotificationDTO struct {
	// Notification identificator
	ID string `json:"id"`
	// Receiver informations
	Receiver UserDTO `json:"receiver"`
	// Timestamp of creation time in seconds (unix)
	CreatedAt int64 `json:"createdAt"`
	// Type of notification system either invitation
	Type enums.NotificationType `json:"type"`
	// Title
	Title string `json:"title"`
	// Description
	Body string `json:"body"`
	// A part of link used to navigate user after clicking notification
	RedirectionLink string `json:"redirectionLink" example:"/meeting/some_id"`
	//Additional information which is unique for other notification types
	Metadata map[metadata.Key]any `json:"metadata"`
}

func (dto *NotificationDTO) JSON() []byte {
	data, err := json.Marshal(*dto)

	if err != nil {
		return []byte{}
	}

	return data
}
