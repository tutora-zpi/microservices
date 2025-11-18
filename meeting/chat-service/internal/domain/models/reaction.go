package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Reaction struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    string        `bson:"userId"`
	MessageID bson.ObjectID `bson:"messageId"`
	SentAt    int64         `bson:"sentAt"`

	Emoji string `bson:"emoji"`
}

func NewReaction(userID, messageID, emoji string, sentAt int64) (*Reaction, error) {
	msgID, err := bson.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, err
	}

	return &Reaction{
		UserID:    userID,
		MessageID: msgID,
		Emoji:     emoji,
		SentAt:    sentAt,
	}, nil
}
