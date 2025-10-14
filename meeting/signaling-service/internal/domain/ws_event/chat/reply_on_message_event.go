package chat

import "github.com/go-playground/validator/v10"

type ReplyOnMessageEvent struct {
	ReplyToMessageID string `json:"replyToMessageID" validate:"reqiured,uuid4"`

	SendMessageEvent
}

func (r *ReplyOnMessageEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(r)
}

func (r *ReplyOnMessageEvent) Type() string {
	return "reply"
}

func (u *ReplyOnMessageEvent) Name() string {
	return u.Type()
}
