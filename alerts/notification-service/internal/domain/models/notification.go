package models

import (
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/metadata"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt int64         `bson:"createdAt"`

	Type   enums.NotificationType   `bson:"type"`
	Status enums.NotificationStatus `bson:"status"`

	Receiver User `bson:"receiver"`
	Sender   User `bson:"sender"`

	Title string `bson:"title"`
	Body  string `bson:"body"`

	RedirectionLink string `bson:"redirectionLink"`

	Metadata map[metadata.Key]any `bson:"metadata"`
}

func NewPartialNotification(notificationType enums.NotificationType, receiverID, senderID string, metadata map[metadata.Key]any) Notification {
	return Notification{
		ID:        bson.NewObjectID(),
		CreatedAt: bson.NewObjectID().Timestamp().Unix(),
		Status:    enums.PENDING,

		Type:     notificationType,
		Receiver: *NewPartialUser(receiverID),
		Sender:   *NewPartialUser(senderID),

		RedirectionLink: "",
		Metadata:        metadata,
	}
}

func (n *Notification) DTO() *dto.NotificationDTO {
	return &dto.NotificationDTO{
		ID:              n.ID.Hex(),
		Receiver:        n.Receiver.DTO(),
		Sender:          n.Sender.DTO(),
		CreatedAt:       n.CreatedAt,
		Type:            n.Type,
		RedirectionLink: n.RedirectionLink,
		Metadata:        n.Metadata,
	}
}

func (n *Notification) SetTitleAndBody(title, body string) {
	n.Title = title
	n.Body = body
}

func (n *Notification) GetHourAndMinute() (hour, minute int) {
	startTime := time.Unix(n.CreatedAt, 0)
	return startTime.Hour(), startTime.Minute()
}
