package meetinginvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/metadata"
	"notification-serivce/internal/domain/models"
	"reflect"
	"time"
)

type MeetingEndedEvent struct {
	MeetingID    string        `json:"meetingId"`
	ClassID      string        `json:"classId"`
	EndTimestamp int64         `json:"endTimestamp"`
	Members      []dto.UserDTO `json:"members"`
}

func (m *MeetingEndedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}

func (m *MeetingEndedEvent) EndedMeetingNotifications() []*models.Notification {
	now := time.Now().UTC().Unix()
	diff := now - m.EndTimestamp

	notifications := make([]*models.Notification, len(m.Members))

	for i, user := range m.Members {
		n := models.BaseNotification()
		n.Title = "Meeting has been finished"
		n.Body = fmt.Sprintf("Your meeting has ended %ds ago", diff)
		n.RedirectionLink = m.buildLink()
		n.Receiver = models.NewUser(user.ID, user.FirstName, user.LastName)
		n.Metadata = map[metadata.Key]any{
			metadata.MEETING_ID: m.MeetingID,
			metadata.CLASS_ID:   m.ClassID,
		}

		notifications[i] = n
	}

	return notifications
}

func (m *MeetingEndedEvent) buildLink() string {
	return fmt.Sprintf("/room/%s", m.ClassID)
}
