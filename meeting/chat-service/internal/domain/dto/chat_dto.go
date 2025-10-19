package dto

import (
	"chat-service/internal/domain/models"
	"time"
)

type ChatDTO struct {
	ID        string       `json:"id"`
	MemberIDs []string     `json:"memberIds"`
	Messages  []MessageDTO `json:"messages"`
	CreatedAt *time.Time   `json:"createdAt,omitempty"`
}

func NewChatDTO(chat models.Chat, messages []MessageDTO) ChatDTO {
	createdAtTime := time.Unix(chat.CreatedAt, 0).UTC()

	return ChatDTO{
		ID:        chat.ID,
		CreatedAt: &createdAtTime,
		MemberIDs: chat.MemberIDs,
		Messages:  messages,
	}
}
