package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	Invitation NotificationType = "invitation"
	System     NotificationType = "system"
)

type NotificationStatus string

const (
	Sent    NotificationStatus = "sent"
	Created NotificationStatus = "created"
)

type Notification struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Type      NotificationType
	Status    NotificationStatus

	ReceiverID string

	Title string
	Body  string

	RedirectionLink string

	Metadata map[string]any
}

func NewNotification(notificationType NotificationType, receiverID, title, body, redirectionLink string, metadata map[string]any) *Notification {
	return &Notification{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Status:    Created,

		Type:       notificationType,
		ReceiverID: receiverID,

		Title:           title,
		Body:            body,
		RedirectionLink: redirectionLink,

		Metadata: metadata,
	}
}
