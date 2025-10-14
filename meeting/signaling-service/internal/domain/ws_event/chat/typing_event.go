package chat

import "github.com/go-playground/validator/v10"

type UserTypingEvent struct {
	ChatID        string `json:"chatID" validate:"required,uuid4"`
	UserTyperID   string `json:"userID" validate:"required,uuid4"`
	UserTyperName string `json:"userTyperName" validate:"required"`
	IsTyping      bool   `json:"isTyping"`
}

func (u *UserTypingEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserTypingEvent) Name() string {
	return "user-typing"
}
