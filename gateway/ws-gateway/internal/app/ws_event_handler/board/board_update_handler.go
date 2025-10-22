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
	"ws-gateway/internal/domain/event"
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
	var wsEvent board.BoardUpdateWSEvent
	if err := json.Unmarshal(body, &wsEvent); err != nil {
		return fmt.Errorf("failed to decode %s payload", wsEvent.Name())
	}

	newEvent := event.NewBoardUpdateEvent(wsEvent)

	wrapper := wsevent.SocketEventWrapper{
		Name:    wsEvent.Name(),
		Payload: body,
	}

	go b.hubManager.Emit(wsEvent.MeetingID, wrapper.ToBytes(), func(id string) bool { return true })

	go b.eventBuffer.Add(newEvent, broker.NewExchangeDestination(newEvent, b.exchange))

	go func() {
		if err := b.cacheService.MakeSnapshot(ctx, wrapper, wsEvent.MeetingID); err != nil {
			log.Printf("Snapshot error: %v", err)
		}
	}()

	return nil
}

func NewBoardUpdateHandler(hubManager interfaces.HubManager, eventBuffer bus.EventBuffer, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	ex := os.Getenv(config.BOARD_EXCHANGE)

	return &boardUpdateHandler{hubManager: hubManager, eventBuffer: eventBuffer, exchange: ex, cacheService: cacheService}
}
