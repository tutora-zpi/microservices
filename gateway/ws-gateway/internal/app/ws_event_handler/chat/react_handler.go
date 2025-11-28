package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/config"
	"ws-gateway/internal/domain/broker"
	"ws-gateway/internal/domain/event"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/chat"
	"ws-gateway/internal/infrastructure/bus"
)

type reactHandler struct {
	hubManager   interfaces.HubManager
	eventBuffer  bus.EventBuffer
	exchange     string
	cacheService interfaces.CacheEventService
}

// Handle implements interfaces.EventHandler.
func (r *reactHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var wsEvent chat.ReactOnMessageWSEvent
	if err := json.Unmarshal(body, &wsEvent); err != nil {
		return fmt.Errorf("failed to decode %s payload", wsEvent.Name())
	}

	newEvent := event.NewReactOnMessageEvent(wsEvent)

	wrapper := wsevent.SocketEventWrapper{
		Name:    wsEvent.Name(),
		Payload: body,
	}

	go r.hubManager.Emit(wsEvent.ChatID, wrapper.ToBytes(), func(id string) bool { return true })

	go r.eventBuffer.Add(newEvent, broker.NewExchangeDestination(newEvent, r.exchange))

	go func() {
		if err := r.cacheService.PushRecentEvent(ctx, wrapper, wsEvent.ChatID); err != nil {
			log.Printf("An error occurred during pushing event: %v", err)
		} else {
			log.Printf("Successfully pushed recent event: %s", wrapper.Name)
		}
	}()

	return nil
}

func NewReactHandler(hubManager interfaces.HubManager, eventBuffer bus.EventBuffer, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	ex := os.Getenv(config.CHAT_EXCHANGE)
	return &reactHandler{hubManager: hubManager, eventBuffer: eventBuffer, exchange: ex, cacheService: cacheService}
}
