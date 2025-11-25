package recorder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/recorder"
)

type stopRecordMeetingHandler struct {
	broker       interfaces.Broker
	exchange     string
	hubManager   interfaces.HubManager
	cacheService interfaces.CacheEventService
}

// Handle implements interfaces.EventHandler.
func (s *stopRecordMeetingHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var evt recorder.StopRecordingRequestedWSEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return fmt.Errorf("failed to decode event: [%s]", evt.Name())
	}

	go func() {
		var payload []byte
		payload, _ = wsevent.EncodeSocketEventWrapper(&evt)
		s.hubManager.Emit(evt.RoomID, payload, func(id string) bool { return true })
	}()

	go func() {
		if err := s.cacheService.RemoveMeetingFromPool(ctx, evt.RoomID); err != nil {
			log.Printf("Failed to remove from meeting pool: %v", err)
		}
	}()

	return nil
}

func NewStopRecordMeetingHandler(broker interfaces.Broker, exchange string, hubManager interfaces.HubManager, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	return &stopRecordMeetingHandler{broker: broker, exchange: exchange, hubManager: hubManager, cacheService: cacheService}
}
