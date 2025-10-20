package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/domain/event"
	wsevent "ws-gateway/internal/domain/ws_event"
)

type sendFileMessageHandler struct {
	hubManager   interfaces.HubManager
	cacheService interfaces.CacheEventService
}

// Handle implements interfaces.EventHandler.
func (s *sendFileMessageHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event event.SendFileMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	wrapper := wsevent.SocketEventWrapper{
		Name:    event.Name(),
		Payload: body,
	}

	go s.hubManager.Emit(event.ChatID, wrapper.ToBytes(), func(id string) bool { return true })

	go func() {
		if err := s.cacheService.PushRecentEvent(ctx, wrapper, event.ChatID); err != nil {
			log.Printf("Snapshot error: %v", err)
		}
	}()

	return nil
}

func NewSendFileMessageHandler(hubManager interfaces.HubManager, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	return &sendFileMessageHandler{hubManager: hubManager, cacheService: cacheService}
}
