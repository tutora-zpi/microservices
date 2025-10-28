package service

import (
	"context"
	"recorder-service/internal/domain/bot"
	"recorder-service/internal/domain/event"
)

type BotService interface {
	AddNewBot(ctx context.Context, evt event.RecordMeetingEvent) (bot.Bot, error)
	GetBot(roomID string) (bot.Bot, bool)
	RemoveBot(ctx context.Context, roomID string) error
}
