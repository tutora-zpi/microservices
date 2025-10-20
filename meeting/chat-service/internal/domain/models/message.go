package models

import (
	"github.com/google/uuid"
)

type Message struct {
	ID     string `bson:"_id,omitempty"`
	SentAt int64  `bson:"sentAt"`

	SenderID string `bson:"senderId"`
	ChatID   string `bson:"chatId"`

	ReplyToID *string `bson:"replyToId,omitempty"`

	ReactionIDs []string `bson:"reactionIds,omitempty"`

	Content string `bson:"content"`

	FileLink *string `bson:"fileLink,omitempty"`
}

func NewMessage(chatID, senderID, content string, sentAt int64, fileLink string) *Message {
	id := uuid.New()

	msg := &Message{
		ID:        id.String(),
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
