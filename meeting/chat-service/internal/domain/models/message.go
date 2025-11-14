package models

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Message struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	SentAt int64         `bson:"sentAt"`

	SenderID string `bson:"senderId"`
	ChatID   string `bson:"chatId"`

	ReplyToID *bson.ObjectID `bson:"replyToId,omitempty"`

	ReactionIDs []bson.ObjectID `bson:"reactionIds,omitempty"`

	Content string `bson:"content"`

	FileLink *string `bson:"fileLink,omitempty"`
}

func NewMessage(chatID, senderID, content, messageID string, sentAt int64, fileLink string) *Message {
	msgID, err := bson.ObjectIDFromHex(messageID)
	if err != nil {
		log.Printf("Failed to cast message id to object id: %v", err)
		return nil
	}

	msg := &Message{
		ID:        msgID,
		SenderID:  senderID,
		ChatID:    chatID,
		SentAt:    sentAt,
		Content:   content,
		ReplyToID: nil,
	}

	if fileLink != "" {
		msg.FileLink = &fileLink
	}

	return msg
}
