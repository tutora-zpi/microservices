package dto

import "chat-service/internal/domain/models"

type ReactionDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	MessageID string `json:"messageId"`
	Emoji     string `json:"emoji"`
}

func NewReactionDTO(reaction models.Reaction) *ReactionDTO {
	return &ReactionDTO{
		ID:        reaction.ID.Hex(),
		UserID:    reaction.UserID,
		MessageID: reaction.MessageID.Hex(),
		Emoji:     reaction.Emoji,
	}
}

func NewReactionDTOs(reactions []models.Reaction) []ReactionDTO {
	var result []ReactionDTO = make([]ReactionDTO, len(reactions))
	for i, reaction := range reactions {
		result[i] = *NewReactionDTO(reaction)
	}

	return result
}
