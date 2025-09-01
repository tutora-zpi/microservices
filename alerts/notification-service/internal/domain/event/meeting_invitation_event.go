package event

import "notification-serivce/internal/domain/models"

type MeetingInvitationEvent struct {
}

func (c *MeetingInvitationEvent) Notification() *models.Notification {
	return &models.Notification{}
}
