package event

import (
	"notification-serivce/internal/domain/models"
	"reflect"
)

type MeetingInvitationEvent struct {
}

func (m *MeetingInvitationEvent) Notification() *models.Notification {
	return &models.Notification{}
}

func (m *MeetingInvitationEvent) Name() string {
	return reflect.TypeOf(m).Elem().Name()
}
