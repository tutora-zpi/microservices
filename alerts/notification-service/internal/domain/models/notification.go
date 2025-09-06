package models

import (
	"log"
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

	Receiver User `bson:"receiver"`
	Sender   User `bson:"sender"`

	RedirectionLink string `bson:"redirection_link"`

	Metadata map[string]any `bson:"metadata"`
}

func NewPartialNotification(notificationType enums.NotificationType, receiverID, senderID string, metadata map[string]any) *Notification {
	return &Notification{
		ID:        bson.NewObjectID(),
		CreatedAt: bson.NewObjectID().Timestamp(),
		Status:    enums.Pending,

		Type:     notificationType,
		Receiver: *NewPartialUser(receiverID),
		Sender:   *NewPartialUser(senderID),

		RedirectionLink: "",
		Metadata:        metadata,
	}
}

func NewNotification() *Notification {
	log.Panic("unimplemented")
	return nil
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
