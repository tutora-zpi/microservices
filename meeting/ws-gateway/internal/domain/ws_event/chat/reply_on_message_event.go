package chat

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ReplyOnMessageEvent struct {
	SendMessageEvent

	ReplyToMessageID string `json:"replyToMessageId" validate:"reqiured,uuid4"`
}

func (r *ReplyOnMessageEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(r)
}

func (r *ReplyOnMessageEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
