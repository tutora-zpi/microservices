package chat

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ReactOnMessageEvent struct {
	MessageID string `json:"messageId" validate:"required,uuid"`

	UserID string `json:"userId" validate:"required,uuid"`

	// Chat ID is create basing on meeting id
	ChatID string `json:"chatId" validate:"required,uuid"`

	Emoji string `json:"emoji" validate:"required,emoji"`

	SentAt int64 `json:"sentAt"`
}

func (r *ReactOnMessageEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(r)
}

func (r *ReactOnMessageEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
