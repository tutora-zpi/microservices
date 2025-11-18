package event

import "reflect"

type SendFileMessageEvent struct {
	SendMessageEvent
}

func NewSendFileMessageEvent(s *SendMessageEvent) *SendFileMessageEvent {
	return &SendFileMessageEvent{
		SendMessageEvent: *s,
	}
}

func (s *SendFileMessageEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
