package classinvitation

import (
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/metadata"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type ClassInvitationCreatedEvent struct {
	ClassName  string `json:"className"`
	ReceiverID string `json:"receiverId"`
	SenderID   string `json:"senderId"`
}

func (c *ClassInvitationCreatedEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}

func (c *ClassInvitationCreatedEvent) Notification() *models.Notification {
	metadata := map[metadata.Key]any{
		metadata.CLASS_NAME: c.ClassName,
	}

	return models.NewPartialNotification(enums.INVITATION, c.ReceiverID, c.SenderID, metadata)
}
