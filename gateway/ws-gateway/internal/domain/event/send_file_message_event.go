package event

import "reflect"

type SendFileMessageEvent struct {
	SendMessageEvent
	FileLink string `json:"fileLink,omitempty"`
	FileName string `json:"fileName,omitempty"`
}

func (s *SendFileMessageEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
