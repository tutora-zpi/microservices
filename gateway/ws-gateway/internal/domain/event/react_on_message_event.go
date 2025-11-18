package event

import (
	"reflect"
	"ws-gateway/internal/domain/ws_event/chat"
)

type ReactOnMessageEvent struct {
	chat.ReactOnMessageWSEvent
}

func NewReactOnMessageEvent(wsevent chat.ReactOnMessageWSEvent) *ReactOnMessageEvent {
	return &ReactOnMessageEvent{
		ReactOnMessageWSEvent: wsevent,
	}
}

func (r *ReactOnMessageEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
