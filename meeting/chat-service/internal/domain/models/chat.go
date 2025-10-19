package models

import (
	"time"
)

type Chat struct {
	ID string `bson:"_id,omitempty"`

	CreatedAt int64 `bson:"createdAt"`

	MemberIDs  []string `bson:"memberIds"`
	MessageIDs []string `bson:"messageIds,omitempty"`
}

func NewChat(memberIDs []string, chatID string) *Chat {
	return &Chat{
		ID:        chatID,
		CreatedAt: time.Now().UTC().Unix(),
		MemberIDs: memberIDs,
	}
}
