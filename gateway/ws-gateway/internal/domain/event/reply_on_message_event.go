package event

import (
	"reflect"
	"ws-gateway/internal/domain/ws_event/chat"
)

type ReplyOnMessageEvent struct {
	chat.ReplyOnMessageWSEvent
}

func NewReplyOnMessageEvent(wsevent chat.ReplyOnMessageWSEvent) *ReplyOnMessageEvent {
	return &ReplyOnMessageEvent{ReplyOnMessageWSEvent: wsevent}
}

func (r *ReplyOnMessageEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
