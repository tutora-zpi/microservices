package event

import "notification-serivce/internal/domain/models"

type ClassInvitationEvent struct {
}

func (c *ClassInvitationEvent) Notification() *models.Notification {
	return &models.Notification{}
}
