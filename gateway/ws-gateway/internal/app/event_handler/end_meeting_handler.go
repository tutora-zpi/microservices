package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/domain/event"
	wsevent "ws-gateway/internal/domain/ws_event"
	recorderDomain "ws-gateway/internal/domain/ws_event/recorder"
)

type endMeetingHandler struct {
	hubManager   interfaces.HubManager
	cacheService interfaces.CacheEventService
}

// Handle implements interfaces.EventHandler.
func (e *endMeetingHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event event.MeetingEndedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	log.Println("Successfully unmarshaled event")

	go func() {
		log.Print("Removing informations about meeting")
		if err := e.cacheService.RemoveMeetingFromPool(ctx, event.MeetingID); err != nil {
			log.Printf("Failed to remove meeting from pool: %v", err)
		}
	}()

	go func() {
		log.Print("Deleting board snapshot")
		if err := e.cacheService.DeleteSnapshot(ctx, event.MeetingID); err != nil {
			log.Printf("Failed to delete snapshot: %v", err)
		}
	}()

	go func() {
		// if bot exist in conn pool he will receive it and make logic if not then not
		stopTime := time.Unix(event.EndTimestamp, 0)

		evt := recorderDomain.StopRecordingRequestedWSEvent{
			RoomID:   event.MeetingID,
			StopTime: stopTime,
		}

		var payload []byte
		payload, _ = wsevent.EncodeSocketEventWrapper(&evt)
		e.hubManager.Emit(event.MeetingID, payload, func(id string) bool { return true })
	}()

	log.Print("Disconnecting users from room")

	for _, user := range event.Members {
		go func() {
			log.Printf("Emitting auto disconnection for %s", user.ID)

			e.hubManager.RemoveRoomMemberByID(event.MeetingID, user.ID)
		}()
	}

	return nil
}

func NewEndMeetingHandler(
	hubManager interfaces.HubManager,
	cacheService interfaces.CacheEventService,
) interfaces.EventHandler {
	return &endMeetingHandler{hubManager: hubManager, cacheService: cacheService}
}
