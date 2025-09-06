package classinvitation

import (
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type ClassInvitationCreatedEvent struct {
	RoomID     string `json:"room_id"`
	RoomName   string `json:"room_name"`
	ReceiverID string `json:"receiver_id"`
	SenderID   string `json:"sender_id"`
}

func (c *ClassInvitationCreatedEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}

func (c *ClassInvitationCreatedEvent) Notification() *models.Notification {
	metadata := map[string]any{
		"room_id":   c.RoomID,
		"room_name": c.RoomName,
	}

	return models.NewPartialNotification(enums.Invitation, c.ReceiverID, c.SenderID, metadata)
}
