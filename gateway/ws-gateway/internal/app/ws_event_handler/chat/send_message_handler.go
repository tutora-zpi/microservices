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

type sendMessageHandler struct {
	hubManager   interfaces.HubManager
	eventBuffer  bus.EventBuffer
	cacheService interfaces.CacheEventService

	exchange string
}

// Handle implements interfaces.EventHandler.
func (s *sendMessageHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var wsEvent chat.SendMessageWSEvent

	if err := json.Unmarshal(body, &wsEvent); err != nil {
		return fmt.Errorf("failed to decode %s payload", wsEvent.Name())
	}

	wsEvent.AppendID()

	newEvent := event.NewSendMessageEvent(wsEvent)

	body, _ = json.Marshal(&wsEvent)

	wrapper := wsevent.SocketEventWrapper{
		Name:    wsEvent.Name(),
		Payload: body,
	}

	go s.hubManager.Emit(wsEvent.ChatID, wrapper.ToBytes(), func(id string) bool { return true })

	go s.eventBuffer.Add(newEvent, broker.NewExchangeDestination(newEvent, s.exchange))

	go func() {
		if err := s.cacheService.PushRecentEvent(ctx, wrapper, wsEvent.ChatID); err != nil {
			log.Printf("An error occurred during pushing event: %v", err)
		} else {
			log.Printf("Successfully pushed recent event: %s", wrapper.Name)
		}
	}()

	return nil
}

func NewSendMessageHandler(
	hubManager interfaces.HubManager,
	eventBuffer bus.EventBuffer,
	cacheService interfaces.CacheEventService,
) interfaces.EventHandler {
	ex := os.Getenv(config.CHAT_EXCHANGE)

	return &sendMessageHandler{
		hubManager:   hubManager,
		eventBuffer:  eventBuffer,
		exchange:     ex,
		cacheService: cacheService,
	}
}
