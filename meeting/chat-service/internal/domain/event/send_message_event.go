package event

import "reflect"

type SendMessageEvent struct {
	MessageID string `json:"messageId"`
	Content   string `json:"content"`
	SenderID  string `json:"senderId"`
	ChatID    string `json:"chatId"`
	SentAt    int64  `json:"sentAt"`
	FileLink  string `json:"fileLink,omitempty"`
	FileName  string `json:"fileName,omitempty"`
}

func (s *SendMessageEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
