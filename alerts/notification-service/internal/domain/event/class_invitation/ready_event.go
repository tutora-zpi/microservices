package classinvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/metadata"
	"reflect"
)

// Domain event
type ClassInvitationReadyEvent struct {
	*dto.NotificationDTO
}

const PLACEHOLDER_ROOM_NAME = "Awesome Room"

func NewClassInvitationReadyEvent(dto *dto.NotificationDTO) *ClassInvitationReadyEvent {
	return &ClassInvitationReadyEvent{
		NotificationDTO: dto,
	}
}

func (c *ClassInvitationReadyEvent) enrichNotification() *dto.NotificationDTO {
	title := fmt.Sprintf("Invitation to %s class!", c.getClassName())

	body := fmt.Sprintf(
		"%s!, You've been invited by %s to %s class. Click button below to go on the invitations page.",
		c.Receiver.FirstName, c.Sender.FullName(), c.getClassName(),
	)

	c.Title = title
	c.Body = body

	return c.NotificationDTO
}

func (c *ClassInvitationReadyEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}

func (c *ClassInvitationReadyEvent) getClassName() string {
	roomName, ok := c.Metadata[metadata.CLASS_NAME].(string)
	if !ok {
		roomName = PLACEHOLDER_ROOM_NAME
	}

	return roomName
}

func (c *ClassInvitationReadyEvent) GetReadyNotification() *dto.NotificationDTO {
	c.enrichNotification()
	return c.NotificationDTO
}
