package classinvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type ClassInvitationCreatedEvent struct {
	ClassName string      `json:"className"`
	Receiver  dto.UserDTO `json:"receiver"`
	Sender    dto.UserDTO `json:"sender"`
}

func (c *ClassInvitationCreatedEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}

func (c *ClassInvitationCreatedEvent) NotificationForReceiver() *models.Notification {
	title := "Got new invivation!"
	body := fmt.Sprintf("%s, you've been invited to %s. Click down below to check new invitations.", c.Receiver.FirstName, c.ClassName)
	link := c.buildLink()

	notification := models.InitInvitationNotification()
	c.buildNotification(notification, title, body, link, c.Receiver)

	return notification
}

func (c *ClassInvitationCreatedEvent) NotificationForSender() *models.Notification {
	title := "Got it!"
	body := fmt.Sprintf("%s, your invitation to %s has been sent.", c.Sender.FirstName, c.ClassName)
	link := c.buildLink()

	notification := models.InitInvitationNotification()
	c.buildNotification(notification, title, body, link, c.Sender)

	return notification
}

func (c *ClassInvitationCreatedEvent) buildNotification(dest *models.Notification, title, body, link string, user dto.UserDTO) {
	dest.Title = title
	dest.Body = body
	dest.RedirectionLink = link
	dest.Receiver = models.NewUser(user.ID, user.FirstName, user.LastName)
}

func (c *ClassInvitationCreatedEvent) buildLink() string {
	return "/dashboard/invitations"
}
