package bot

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type Bot struct {
	ID   string `json:"botId"`
	Name string `json:"botName" fake:"{firstname}"`
}

func NewBot() *Bot {
	var b Bot
	err := gofakeit.Struct(&b)
	if err != nil {
		return nil
	}

	b.ID = uuid.NewString()
	return &b
}
