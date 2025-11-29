package meetinginvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/models"
	"reflect"
	"time"
)

type MeetingEndedEvent struct {
	MeetingID    string        `json:"meetingId"`
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
		n.Body = fmt.Sprintf("%s!, your meeting has ended %ds ago", user.FirstName, diff)
		n.RedirectionLink = m.buildLink()
		n.Receiver = models.NewUser(user.ID, user.FirstName, user.LastName)

		notifications[i] = n
	}

	return notifications
}

func (m *MeetingEndedEvent) buildLink() string {
	return ""
}
