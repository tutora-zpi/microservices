package board

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
	"signaling-service/internal/domain/ws_event/board"
	"signaling-service/internal/infrastructure/bus"
)

type boardUpdateHandler struct {
	hubManager   interfaces.HubManager
	eventBuffer  bus.EventBuffer
	cacheService interfaces.CacheEventService

	exchange string
}

// Handle implements interfaces.EventHandler.
func (b *boardUpdateHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event board.BoardUpdateEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	go b.hubManager.Emit(event.MeetingID, body, func(id string) bool { return true })

	b.eventBuffer.Add(&event, broker.NewExchangeDestination(&event, b.exchange))

	wrapper := wsevent.SocketEventWrapper{
		Name:    event.Name(),
		Payload: body,
	}

	go func() {
		if err := b.cacheService.MakeSnapshot(ctx, wrapper, event.MeetingID); err != nil {
			log.Printf("Snapshot error: %v", err)
		}
	}()

	return nil
}

func NewBoardUpdateHandler(hubManager interfaces.HubManager, eventBuffer bus.EventBuffer, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	ex := os.Getenv(config.BOARD_EXCHANGE)

	return &boardUpdateHandler{hubManager: hubManager, eventBuffer: eventBuffer, exchange: ex, cacheService: cacheService}
}
