package event

import "reflect"

type ReplyOnMessageEvent struct {
	SendMessageEvent

	ReplyToMessageID string `json:"replyToMessageId"`
}

func (r *ReplyOnMessageEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
