package models

import (
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/enums"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"createdAt"`

	Type   enums.NotificationType   `bson:"type"`
	Status enums.NotificationStatus `bson:"status"`

	ReceiverID string `bson:"receiverId"`

	Title           string `bson:"title"`
	Body            string `bson:"body"`
	RedirectionLink string `bson:"redirectionLink"`

	Metadata map[string]any `bson:"metadata"`
}

func NewNotification(notificationType enums.NotificationType, receiverID, title, body, redirectionLink string, metadata map[string]any) *Notification {
	return &Notification{
		ID:        bson.NewObjectID(),
		CreatedAt: bson.NewObjectID().Timestamp(),
		Status:    enums.Created,

		Type:       notificationType,
		ReceiverID: receiverID,

		Title:           title,
		Body:            body,
		RedirectionLink: redirectionLink,

		Metadata: metadata,
	}
}

func (n *Notification) DTO() dto.NotificationDTO {
	return dto.NotificationDTO{
		ID:              n.ID.Hex(),
		CreatedAt:       n.CreatedAt,
		Type:            n.Type,
		Status:          n.Status,
		Title:           n.Title,
		Body:            n.Body,
		RedirectionLink: n.RedirectionLink,
		Metadata:        n.Metadata,
	}
}
