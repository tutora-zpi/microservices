package bot

import (
	"recorder-service/internal/domain/client"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/recorder"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type Bot interface {
	ID() string
	Name() string

	DTO() dto.BotDTO

	Client() client.Client
	Recorder() recorder.Recorder
}

type bot struct {
	id   string
	name string `fake:"{firstname}"`

	recorder recorder.Recorder
	client   client.Client
}

// DTO implements Bot.
func (b bot) DTO() dto.BotDTO {
	return dto.BotDTO{
		ID:   b.id,
		Name: b.name,
	}
}

// Client implements Bot.
func (b bot) Client() client.Client {
	return b.client
}

// ID implements Bot.
func (b bot) ID() string {
	return b.id
}

// Name implements Bot.
func (b bot) Name() string {
	return b.name
}

// Recorder implements Bot.
func (b bot) Recorder() recorder.Recorder {
	return b.recorder
}

func NewBot(rec recorder.Recorder, client client.Client) Bot {
	var b bot

	b.recorder = rec
	b.client = client

	err := gofakeit.Struct(&b)
	if err != nil {
		return nil
	}

	b.id = uuid.NewString()
	return b
}
