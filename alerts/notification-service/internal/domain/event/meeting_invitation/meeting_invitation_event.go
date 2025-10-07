package meetinginvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/models"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MeetingStartedEvent struct {
	MeetingID   string        `json:"meetingId"`
	Members     []dto.UserDTO `json:"members"`
	StartedTime time.Time     `json:"startedTime"` // ISO 8601 format
}

func (m *MeetingStartedEvent) Notifications() []models.Notification {
	base := models.Notification{
		ID:              bson.NewObjectID(),
		CreatedAt:       m.StartedTime.Unix(),
		Type:            enums.INVITATION,
		Status:          enums.CREATED,
		RedirectionLink: m.buildLink(),
		Metadata:        nil,
		Sender:          models.User{},
	}

	notifications := []models.Notification{}

	for _, receiver := range m.Members {
		notification := base
		notification.ID = bson.NewObjectID()
		hour, minute := notification.GetHourAndMinute()
		notification.Receiver = *models.NewUser(receiver.ID, receiver.FirstName, receiver.LastName, "")

		title := fmt.Sprintf("%s, meeting has already started!", notification.Receiver.FirstName)
		body := fmt.Sprintf("Meeting was scheduled on %02d:%02d. Click down below to join!",
			hour, minute)

		notification.SetTitleAndBody(title, body)

		notifications = append(notifications, notification)
	}

	return notifications
}

func (m *MeetingStartedEvent) Name() string {
	return reflect.TypeOf(m).Elem().Name()
}

func (m *MeetingStartedEvent) buildLink() string {
	return fmt.Sprintf("/meeting/%s", m.MeetingID)
}
