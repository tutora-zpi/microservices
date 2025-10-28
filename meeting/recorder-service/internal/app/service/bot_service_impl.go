package service

import (
	"context"
	"recorder-service/internal/app/interfaces/factory"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/bot"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/repository"
	"sync"
	"time"
)

type botServiceImpl struct {
	rooms map[string]bot.Bot
	mutex sync.Mutex

	botCache repository.BotRepository

	recorderFactory factory.RecorderFactory
	clientFactory   factory.ClientFactory
}

// AddNewBot implements interfaces.BotService.
func (b *botServiceImpl) AddNewBot(ctx context.Context, evt event.RecordMeetingEvent) (bot.Bot, error) {
	recorder := b.recorderFactory.CreateNewRecorder()
	client := b.clientFactory.CreateNewClient()

	bot := bot.NewBot(recorder, client)

	ttl := time.Until(evt.FinishTime)
	err := b.botCache.TryAdd(ctx, evt.RoomID, bot.ID(), ttl)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

// RemoveBot implements interfaces.BotService.
func (b *botServiceImpl) RemoveBot(ctx context.Context, roomID string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	delete(b.rooms, roomID)
	return b.botCache.Delete(ctx, roomID)
}

// GetBot implements interfaces.BotService.
func (b *botServiceImpl) GetBot(roomID string) (bot.Bot, bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	bt, ok := b.rooms[roomID]
	return bt, ok
}

func NewBotService(
	repo repository.BotRepository,
	recorderFactory factory.RecorderFactory,
	clientFactory factory.ClientFactory,
) service.BotService {
	return &botServiceImpl{
		botCache:        repo,
		rooms:           make(map[string]bot.Bot),
		clientFactory:   clientFactory,
		recorderFactory: recorderFactory,
	}
}
