package board

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/config"
	"ws-gateway/internal/domain/broker"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/board"
	"ws-gateway/internal/infrastructure/bus"
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

	wrapper := wsevent.SocketEventWrapper{
		Name:    event.Name(),
		Payload: body,
	}

	go b.hubManager.Emit(event.MeetingID, wrapper.ToBytes(), func(id string) bool { return true })

	b.eventBuffer.Add(&event, broker.NewExchangeDestination(&event, b.exchange))

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
