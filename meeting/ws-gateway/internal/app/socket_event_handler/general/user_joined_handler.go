package general

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/general"
)

type userJoinedHandler struct {
	hubManager   interfaces.HubManager
	cacheService interfaces.CacheEventService
}

// Handle implements interfaces.EventHandler.
func (u *userJoinedHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event general.UserJoinedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	ids := u.hubManager.AddRoomMember(event.RoomID, client)

	roomUsers := general.RoomUsersEvent{Users: ids}
	bytes, _ := wsevent.EncodeSocketEventWrapper(&roomUsers, roomUsers.Name())

	go u.hubManager.Emit(event.RoomID, bytes, func(id string) bool { return true })

	go func() {
		payloads, eventsErr := u.cacheService.GetLastEventsData(ctx, event.RoomID)
		snapshot, snapErr := u.cacheService.GetSnapshot(ctx, event.RoomID)

		if snapErr != nil || eventsErr != nil {
			log.Printf("Failed to flush data: %v, %v", snapErr, eventsErr)
			return
		}

		payloads = append(payloads, snapshot)
		u.hubManager.EmitToClientInRoom(event.RoomID, client.ID(), payloads)
	}()

	return nil
}

func NewUserJoinedHandler(hubManager interfaces.HubManager, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	return &userJoinedHandler{hubManager: hubManager, cacheService: cacheService}
}
