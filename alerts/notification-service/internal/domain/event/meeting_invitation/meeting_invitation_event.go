package meetinginvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/metadata"
	"notification-serivce/internal/domain/models"
	"reflect"
	"time"
)

type MeetingStartedEvent struct {
	ClassID     string        `json:"classId"`
	MeetingID   string        `json:"meetingId"`
	Members     []dto.UserDTO `json:"members"`
	StartedTime time.Time     `json:"startedTime"` // ISO 8601 format
	FinishTime  time.Time     `json:"finishTime"`
}

func (m *MeetingStartedEvent) Notifications() []models.Notification {
	notifications := make([]models.Notification, len(m.Members))

	for i, receiver := range m.Members {
		notifications[i] = *m.MeetingNotification(receiver)
	}

	return notifications
}

func (m *MeetingStartedEvent) MeetingNotification(user dto.UserDTO) *models.Notification {
	base := models.InitInvitationNotification()
	base.Receiver = models.NewUser(user.ID, user.FirstName, user.LastName)

	base.Title = fmt.Sprintf("%s, meeting has already started!", user.FirstName)

	hour := m.StartedTime.Hour()
	minute := m.StartedTime.Minute()
	duration := m.FinishTime.Sub(m.StartedTime).String()

	base.Body = fmt.Sprintf("Meeting was scheduled on %02d:%02d and duration time is %s. Click down below to join!",
		hour, minute, duration)

	base.RedirectionLink = m.buildLink()

	base.Metadata[metadata.CLASS_ID] = m.ClassID

	return base
}

func (m *MeetingStartedEvent) Name() string {
	return reflect.TypeOf(m).Elem().Name()
}

func (m *MeetingStartedEvent) buildLink() string {
	return fmt.Sprintf("/room/%s/meeting/%s", m.ClassID, m.MeetingID)
}
