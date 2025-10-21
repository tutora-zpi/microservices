package event

import (
	"reflect"
	"ws-gateway/internal/domain/ws_event/chat"
)

type SendMessageEvent struct {
	chat.SendMessageWSEvent
}

func NewSendMessageEvent(wsevent chat.SendMessageWSEvent) *SendMessageEvent {
	return &SendMessageEvent{SendMessageWSEvent: wsevent}
}

func (s *SendMessageEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
