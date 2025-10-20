package chat

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SendMessageWSEvent struct {
	Content   string `json:"content" validate:"required"`
	SenderID  string `json:"senderId" validate:"required,uuid4"`
	ChatID    string `json:"chatId" validate:"required,uuid4"`
	SentAt    int64  `json:"sentAt"`
	MessageID string `json:"messageId,omitempty"`
}

func (s *SendMessageWSEvent) AppendID() {
	s.MessageID = uuid.NewString()
}

func (s *SendMessageWSEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(s)
}

func (s *SendMessageWSEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
