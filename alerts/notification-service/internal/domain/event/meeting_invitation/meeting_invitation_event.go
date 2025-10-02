package meetinginvitation

import (
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
		RedirectionLink: "",
		Metadata:        nil,
		Sender:          models.User{},
	}

	notifications := []models.Notification{}

	for _, receiver := range m.Members {
		notification := base
		notification.ID = bson.NewObjectID()
		notification.Receiver = *models.NewUser(receiver.ID, receiver.FirstName, receiver.LastName, "")
		notifications = append(notifications, notification)
	}

	return notifications
}

func (m *MeetingStartedEvent) Name() string {
	return reflect.TypeOf(m).Elem().Name()
}
