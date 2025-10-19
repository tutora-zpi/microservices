package models

import (
	"github.com/google/uuid"
)

type Reaction struct {
	ID        string `bson:"_id,omitempty"`
	UserID    string `bson:"userId"`
	MessageID string `bson:"messageId"`
	SentAt    int64  `bson:"sentAt"`

	Emoji string `bson:"emoji"`
}

func NewReaction(userID, messageID, emoji string, sentAt int64) (*Reaction, error) {
	id := uuid.New().String()

	return &Reaction{
		ID:        id,
		UserID:    userID,
		MessageID: messageID,
		Emoji:     emoji,
		SentAt:    sentAt,
	}, nil
}
