package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/config"
	"signaling-service/internal/domain/broker"
	wsevent "signaling-service/internal/domain/ws_event"
	"signaling-service/internal/domain/ws_event/chat"
	"signaling-service/internal/infrastructure/bus"
)

type sendMessageHandler struct {
	hubManager   interfaces.HubManager
	eventBuffer  bus.EventBuffer
	cacheService interfaces.CacheEventService

	exchange string
}

// Handle implements interfaces.EventHandler.
func (s *sendMessageHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event chat.SendMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	go s.hubManager.Emit(event.ChatID, body, func(id string) bool { return true })

	s.eventBuffer.Add(&event, broker.NewExchangeDestination(&event, s.exchange))

	wrapper := wsevent.SocketEventWrapper{
		Name:    event.Name(),
		Payload: body,
	}

	go func() {
		if err := s.cacheService.PushRecentEvent(ctx, wrapper, event.ChatID); err != nil {
			log.Printf("An error occurred during pushing event: %v", err)
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
