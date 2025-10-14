package chat

import "github.com/go-playground/validator/v10"

type ReactOnMessageEvent struct {
	MessageID string `json:"messageID" validate:"required,uuid"`

	UserID string `json:"userID" validate:"required,uuid"`

	// Chat ID is create basing on meeting id
	ChatID string `json:"chatID" validate:"required,uuid"`

	Emoji string `json:"emoji" validate:"required,emoji"`
}

func (r *ReactOnMessageEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(r)
}

func (r *ReactOnMessageEvent) Name() string {
	return "react"
}
