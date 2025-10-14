package chat

import "github.com/go-playground/validator/v10"

type SendMessageEvent struct {
	Content   string `json:"content" validate:"required,min=1,max=100"`
	SenderID  string `json:"senderId" validate:"required,uuid4"`
	MeetingID string `json:"meetingId" validate:"required,uuid4"`
}

func (s *SendMessageEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(s)
}

func (s *SendMessageEvent) Type() string {
	return "send-message"
}

func (u *SendMessageEvent) Name() string {
	return u.Type()
}
