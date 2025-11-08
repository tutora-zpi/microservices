package service

import (
	"context"
	"recorder-service/internal/domain/bot"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/recorder"
	"recorder-service/internal/domain/ws_event/general"
)

type BotService interface {
	AddNewBot(ctx context.Context, evt event.RecordMeetingEvent) (bot.Bot, error)
	GetBot(roomID string) (bot.Bot, bool)
	RemoveBot(ctx context.Context, roomID string) ([]recorder.RecordingInfo, error)
	UpdateBotPeers(evt general.RoomUsersWSEvent) error
}
