package chat

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type UserTypingWSEvent struct {
	ChatID        string `json:"chatID" validate:"required,uuid4"`
	UserTyperID   string `json:"userID" validate:"required,uuid4"`
	UserTyperName string `json:"userTyperName" validate:"required"`
	IsTyping      bool   `json:"isTyping"`
}

func (u *UserTypingWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserTypingWSEvent) Name() string {
	return reflect.TypeOf(*u).Name()
}
