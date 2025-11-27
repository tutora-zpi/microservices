package event

import (
	"reflect"
	"ws-gateway/internal/domain/ws_event/chat"
)

type ReplyOnMessageEvent struct {
	chat.ReplyOnMessageWSEvent
}

func NewReplyOnMessageEvent(wsevent chat.ReplyOnMessageWSEvent) *ReplyOnMessageEvent {
	evt := &ReplyOnMessageEvent{ReplyOnMessageWSEvent: wsevent}
	evt.AppendID()
	return evt
}

func (r *ReplyOnMessageEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
