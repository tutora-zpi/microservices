package event

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ReactOnMessageEvent struct {
	MessageID string `json:"messageId"`
	UserID    string `json:"userID"`
	Emoji     string `json:"emoji"`
	ChatID    string `json:"chatId"`
	SentAt    int64  `json:"sentAt"`
}

func (e *ReactOnMessageEvent) IsValid() error {
	v := validator.New()

	return v.Struct(e)
}

func (r *ReactOnMessageEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
