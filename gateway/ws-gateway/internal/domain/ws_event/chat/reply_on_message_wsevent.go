package chat

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ReplyOnMessageWSEvent struct {
	SendMessageWSEvent

	ReplyToMessageID string `json:"replyToMessageId" validate:"reqiured"`
}

func (r *ReplyOnMessageWSEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(r)
}

func (r *ReplyOnMessageWSEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
