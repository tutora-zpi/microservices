package classinvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type ClassInvitationAcceptedEvent struct {
	ClassID   string      `json:"classId"`
	ClassName string      `json:"className"`
	Receiver  dto.UserDTO `json:"accepter"`
	Sender    dto.UserDTO `json:"roomHost"`
}

func (c *ClassInvitationAcceptedEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}

func (c *ClassInvitationAcceptedEvent) NotificationForSender() *models.Notification {
	title := fmt.Sprintf("%s, your invitation has been accepted by %s", c.Sender.FirstName, c.Receiver.FirstName)
	body := fmt.Sprintf("%s accepted your invitation to %s", c.Receiver.FirstName, c.ClassName)

	base := models.BaseNotification()

	base.Title = title
	base.Body = body
	base.RedirectionLink = c.buildLink()

	base.Receiver = models.NewUser(c.Sender.ID, c.Sender.FirstName, c.Sender.LastName)

	return base
}

func (c *ClassInvitationAcceptedEvent) NotificationForReceiver() *models.Notification {
	title := "Successfully accepted invitation"
	body := fmt.Sprintf("Accepted invitation to %s from %s", c.ClassName, c.Sender.FirstName)

	base := models.BaseNotification()

	base.Title = title
	base.Body = body
	base.RedirectionLink = c.buildLink()

	base.Receiver = models.NewUser(c.Receiver.ID, c.Receiver.FirstName, c.Receiver.LastName)

	return base
}

func (c *ClassInvitationAcceptedEvent) buildLink() string {
	return fmt.Sprintf("/room/%s", c.ClassID)
}
