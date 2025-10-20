package requests

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type SaveFileMessage struct {
	Content  string `form:"content"`
	SentAt   int    `form:"sentAt" validate:"gt=0"`
	ChatID   string `form:"chatId" validate:"required,uuid4"`
	SenderID string `form:"senderId" validate:"required,uuid4"`
}

func (s *SaveFileMessage) IsValid() error {
	validator := validator.New()

	return validator.Struct(s)
}

func NewSaveFileMessage(content, senderID, chatID, sentAt string) (*SaveFileMessage, error) {
	v, err := strconv.Atoi(sentAt)
	if err != nil {
		return nil, fmt.Errorf("Invliad sentAt value - NaN")
	}

	return &SaveFileMessage{
		Content:  content,
		SentAt:   v,
		SenderID: senderID,
		ChatID:   chatID,
	}, nil
}
