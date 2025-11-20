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
	log.Print(string(body))
	var event general.UserJoinedWSEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	log.Printf("User with id: %s tries to join to room: %s", event.UserID, event.RoomID)

	ids := u.hubManager.AddRoomMember(event.RoomID, client)

	roomUsers := &general.RoomUsersWSEvent{Users: ids, RoomID: event.RoomID}
	bytes, _ := wsevent.EncodeSocketEventWrapper(roomUsers)

	go u.hubManager.Emit(event.RoomID, bytes, func(id string) bool { return true })

	go func() {
		info, err := u.cacheService.IsMeetingRecorded(ctx, event.RoomID)
		if err != nil {
			log.Printf("Something went during getting meeting recoreded")
			return
		}

		log.Printf("Sending to %s info: %v", client.ID(), *info)

		wrapped, _ := wsevent.EncodeSocketEventWrapper(info)

		u.hubManager.EmitToClientInRoom(event.RoomID, client.ID(), [][]byte{wrapped})
	}()

	go func() {
		payloads, eventsErr := u.cacheService.GetLastEventsData(ctx, event.RoomID)
		snapshot, snapErr := u.cacheService.GetSnapshot(ctx, event.RoomID)

		if snapErr != nil {
			log.Println("Snapshot not found")
		} else {
			payloads = append(payloads, snapshot)
		}

		if eventsErr != nil {
			log.Println("There werent events on server")
		}

		u.hubManager.EmitToClientInRoom(event.RoomID, client.ID(), payloads)
	}()

	return nil
}

func NewUserJoinedHandler(hubManager interfaces.HubManager, cacheService interfaces.CacheEventService) interfaces.EventHandler {
	return &userJoinedHandler{hubManager: hubManager, cacheService: cacheService}
}
