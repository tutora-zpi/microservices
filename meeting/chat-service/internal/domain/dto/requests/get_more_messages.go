package requests

import "strconv"

const DEFAULT_LIMIT = 10

type GetMoreMessages struct {
	ID            string  `json:"id" validate:"required,uuid4"`
	Limit         int     `json:"limit"`
	LastMessageId *string `json:"lastMessageId,omitempty"`
}

func NewGetMoreMessages(id, limit, lastMessageID string) *GetMoreMessages {
	result := &GetMoreMessages{
		ID:    id,
		Limit: DEFAULT_LIMIT,
	}

	if lastMessageID != "" {
		result.LastMessageId = &lastMessageID
	}

	l, err := strconv.Atoi(limit)
	if err == nil {
		result.Limit = l
	}

	return result
}
