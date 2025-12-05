package meetinginvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/metadata"
	"notification-serivce/internal/domain/models"
	"reflect"
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
	notifications := make([]*models.Notification, len(m.Members))

	for i, user := range m.Members {
		n := models.BaseNotification()
		n.Title = "Meeting has been finished"
		n.Body = "You will be redirected to class"
		n.RedirectionLink = m.buildLink()
		n.Receiver = models.NewUser(user.ID, user.FirstName, user.LastName)
		n.Metadata = map[metadata.Key]any{
			metadata.MEETING_ID:    m.MeetingID,
			metadata.CLASS_ID:      m.ClassID,
			metadata.AUTO_REDIRECT: true,
		}

		notifications[i] = n
	}

	return notifications
}

func (m *MeetingEndedEvent) buildLink() string {
	return fmt.Sprintf("/room/%s", m.ClassID)
}
