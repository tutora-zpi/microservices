package event

import (
	"fmt"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type ClassInvitationEvent struct {
	SenderFullName string `json:"sender_full_name"`
	RoomName       string `json:"room_name"`
	ReceiverID     string `json:"receiver_id"`
}

func (c *ClassInvitationEvent) Notification() *models.Notification {
	title := fmt.Sprintf("Invitation to %s class!", c.RoomName)
	body := fmt.Sprintf("You've been invited by %s to %s class. Click button below to go on the invitations page.", c.SenderFullName, c.RoomName)
	link := ""

	return models.NewNotification(enums.Invitation, c.ReceiverID, title, body, link, nil)
}

func (c *ClassInvitationEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}
