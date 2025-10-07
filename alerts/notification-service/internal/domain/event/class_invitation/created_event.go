package classinvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/metadata"
	"notification-serivce/internal/domain/models"
	"reflect"
	"time"
)

type ClassInvitationCreatedEvent struct {
	ClassName string      `json:"className"`
	Receiver  dto.UserDTO `json:"receiver"`
	Sender    dto.UserDTO `json:"sender"`
}

func (c *ClassInvitationCreatedEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}

func (c *ClassInvitationCreatedEvent) Notification() models.Notification {
	title := fmt.Sprintf("%s, meeting has already started!", c.Receiver.FirstName)
	link := c.buildLink()

	notification := models.NewNotification(enums.INVITATION, c.Receiver, c.Sender, title, "", link, map[metadata.Key]any{})
	time.Unix(notification.CreatedAt, 0)

	hour, minute := notification.GetHourAndMinute()

	notification.Body = fmt.Sprintf("Meeting was scheduled on %02d:%02d. Click down below to join!",
		hour, minute)

	return notification
}

func (c *ClassInvitationCreatedEvent) buildLink() string {
	return "/dashboard/invitations"
}
