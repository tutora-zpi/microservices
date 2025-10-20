package event

import "reflect"

type SendMessageEvent struct {
	Content  string `json:"content"`
	SenderID string `json:"senderId"`
	ChatID   string `json:"chatId"`
	SentAt   int64  `json:"sentAt"`
	FileLink string `json:"fileLink,omitempty"`
}

func (s *SendMessageEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
