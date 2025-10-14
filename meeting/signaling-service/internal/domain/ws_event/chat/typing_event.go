package chat

import "github.com/go-playground/validator/v10"

type UserTypingEvent struct {
	ChatID   string `json:"chatID" validate:"required,uuid4"`
	UserID   string `json:"userID" validate:"required,uuid4"`
	IsTyping bool   `json:"isTyping"`
}

func (u *UserTypingEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserTypingEvent) Type() string {
	return "user-typing"
}

func (u *UserTypingEvent) Name() string {
	return u.Type()
}
