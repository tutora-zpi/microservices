package models

import (
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/enums"
	"notification-serivce/internal/domain/metadata"
	"notification-serivce/pkg"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt int64         `bson:"createdAt"`

	Type   enums.NotificationType   `bson:"type"`
	Status enums.NotificationStatus `bson:"status"`

	Receiver User `bson:"receiver"`

	Title string `bson:"title"`
	Body  string `bson:"body"`

	RedirectionLink string `bson:"redirectionLink"`

	Metadata map[metadata.Key]any `bson:"metadata"`
}

func (n *Notification) DTO() *dto.NotificationDTO {
	return &dto.NotificationDTO{
		ID:              n.ID.Hex(),
		Receiver:        n.Receiver.DTO(),
		CreatedAt:       n.CreatedAt,
		Type:            n.Type,
		RedirectionLink: n.RedirectionLink,
		Metadata:        n.Metadata,
		Title:           n.Title,
		Body:            n.Body,
	}
}

func InitInvitationNotification() *Notification {
	return &Notification{
		ID:        bson.NewObjectID(),
		CreatedAt: pkg.GenerateTimestamp(),
		Type:      enums.INVITATION,
		Status:    enums.CREATED,
		Metadata:  map[metadata.Key]any{},
	}
}

func BaseNotification() *Notification {
	return &Notification{
		ID:        bson.NewObjectID(),
		CreatedAt: pkg.GenerateTimestamp(),
		Type:      enums.SYSTEM,
		Status:    enums.CREATED,
		Metadata:  map[metadata.Key]any{},
	}
}
