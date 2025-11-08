package service

import (
	"context"
	"fmt"
	"log"
	"recorder-service/internal/app/interfaces/factory"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/bot"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/recorder"
	"recorder-service/internal/domain/repository"
	"recorder-service/internal/domain/ws_event/general"
	"recorder-service/internal/infrastructure/webrtc/peer"
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
	log.Printf("AddNewBot called for room %s", evt.RoomID)

	recorder := b.recorderFactory.CreateNewRecorder()
	client := b.clientFactory.CreateNewClient()

	newBot := bot.NewBot(client)
	newBot.Client().Connect(ctx)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if existingBot, exists := b.rooms[evt.RoomID]; exists {
		log.Printf("Bot already exists for room %s, returning existing", evt.RoomID)
		return existingBot, nil
	}

	ttl := time.Until(evt.FinishTime)
	if err := b.botCache.TryAdd(ctx, evt.RoomID, newBot.ID(), ttl); err != nil {
		log.Printf("Failed to add bot to cache for room %s: %v", evt.RoomID, err)
	}

	b.rooms[evt.RoomID] = newBot
	log.Printf("Added new bot %s for room %s (TTL: %s)", newBot.ID(), evt.RoomID, ttl)

	for _, id := range evt.ExpectedUserIDs {
		p := peer.NewPeer(evt.RoomID, newBot.ID(), id, recorder, newBot.Client().Send)
		if err := newBot.AddPeer(id, p); err != nil {
			return nil, fmt.Errorf("failed add peer: %v", err)
		}
	}

	return newBot, nil
}

func (b *botServiceImpl) UpdateBotPeers(evt general.RoomUsersWSEvent) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	bot, ok := b.rooms[evt.RoomID]
	if !ok {
		return fmt.Errorf("no bot found in room: %s", evt.RoomID)
	}

	recorder := b.recorderFactory.CreateNewRecorder()

	for _, id := range evt.Users {
		_, ok := bot.GetPeer(id)
		if !ok {
			newPeer := peer.NewPeer(evt.RoomID, bot.ID(), id, recorder, bot.Client().Send)
			if err := bot.AddPeer(id, newPeer); err != nil {
				log.Printf("Failed to add peer connection to bot: %v", err)
			}
		}
	}

	return nil
}

// RemoveBot implements interfaces.BotService.
func (b *botServiceImpl) RemoveBot(ctx context.Context, roomID string) ([]recorder.RecordingInfo, error) {
	log.Printf("Removing bot for room %s", roomID)

	b.mutex.Lock()
	defer b.mutex.Unlock()
	bot, ok := b.rooms[roomID]

	if !ok {
		log.Printf("No bot found for room %s", roomID)
		return nil, nil
	} else {
		delete(b.rooms, roomID)
		log.Printf("Bot removed from memory for room %s", roomID)
	}

	if err := b.botCache.Delete(ctx, roomID); err != nil {
		log.Printf("Failed to delete bot from cache for room %s: %v", roomID, err)
		return nil, err
	}
	log.Printf("Bot cache entry deleted for room %s", roomID)

	infos := bot.FinishRecording()
	if len(infos) < 1 {
		return nil, fmt.Errorf("no info's about recordings in room: %s", roomID)
	}

	return infos, nil
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
	log.Println("Initializing new bot service")
	return &botServiceImpl{
		botCache:        repo,
		rooms:           make(map[string]bot.Bot),
		clientFactory:   clientFactory,
		recorderFactory: recorderFactory,
	}
}
