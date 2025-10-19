package chat

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type SendMessageEvent struct {
	Content  string `json:"content" validate:"required,min=1,max=100"`
	SenderID string `json:"senderId" validate:"required,uuid4"`
	ChatID   string `json:"chatId" validate:"required,uuid4"`
	SentAt   int64  `json:"sentAt"`
}

func (s *SendMessageEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(s)
}

func (s *SendMessageEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
