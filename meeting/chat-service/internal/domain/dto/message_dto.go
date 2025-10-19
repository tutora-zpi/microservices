package dto

import (
	"chat-service/internal/domain/models"
	"time"
)

// MessageDTO represents a message in a chat.
type MessageDTO struct {
	// ID of the message
	// required: true
	ID string `json:"id"`

	// Timestamp when the message was sent
	// required: true
	SentAt time.Time `json:"sentAt"`

	// Content of the message
	// required: false
	Content string `json:"content,omitempty"`

	// ID of the sender
	// required: false
	SenderID string `json:"senderId,omitempty"`

	// ID of the chat
	// required: false
	ChatID string `json:"chatId,omitempty"`

	// Reactions to the message
	// required: false
	Reactions []ReactionDTO `json:"reactions,omitempty"`

	// Optional reply to another message
	// required: false
	Reply *MessageDTO `json:"reply,omitempty"`
}

func NewMessageDTO(message models.Message, replyTo *models.Message, reactions []models.Reaction) *MessageDTO {
	sentAtTime := time.Unix(message.SentAt, 0).UTC()

	reactionDTOs := make([]ReactionDTO, len(reactions))
	for i, r := range reactions {
		reactionDTOs[i] = *NewReactionDTO(r)
	}

	var replyDTO *MessageDTO
	if replyTo != nil {
		replyDTO = &MessageDTO{
			ID:      replyTo.ID,
			Content: replyTo.Content,
		}
	}

	return &MessageDTO{
		ID:        message.ID,
		SentAt:    sentAtTime,
		SenderID:  message.SenderID,
		ChatID:    message.ChatID,
		Content:   message.Content,
		Reply:     replyDTO,
		Reactions: reactionDTOs,
	}
}
