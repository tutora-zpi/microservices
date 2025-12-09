package classinvitation

import (
	"fmt"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type ClassInvitationAcceptedEvent struct {
	ClassID   string `json:"classId"`
	ClassName string `json:"className"`
	Receiver  string `json:"accepterId"`
	Sender    string `json:"roomHostId"`
}

func (c *ClassInvitationAcceptedEvent) Name() string {
	return reflect.TypeOf(c).Elem().Name()
}

func (c *ClassInvitationAcceptedEvent) NotificationForSender() *models.Notification {
	title := "Your invitation has been accepted"
	body := fmt.Sprintf("Accepted your invitation to %s", c.ClassName)

	base := models.BaseNotification()

	base.Title = title
	base.Body = body
	base.RedirectionLink = c.buildLink()

	base.Receiver = models.NewUser(c.Sender, "", "")

	return base
}

func (c *ClassInvitationAcceptedEvent) NotificationForReceiver() *models.Notification {
	title := "Successfully accepted invitation"
	body := fmt.Sprintf("Accepted invitation from %s", c.ClassName)

	base := models.BaseNotification()

	base.Title = title
	base.Body = body
	base.RedirectionLink = c.buildLink()

	base.Receiver = models.NewUser(c.Receiver, "", "")

	return base
}

func (c *ClassInvitationAcceptedEvent) buildLink() string {
	return fmt.Sprintf("/room/%s", c.ClassID)
}
