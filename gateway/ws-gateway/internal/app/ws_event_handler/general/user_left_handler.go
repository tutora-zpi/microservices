package general

import (
	"context"
	"encoding/json"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/general"
)

type userLeftHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *userLeftHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event general.UserLeftWSEvent

	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}

	ids := u.hubManager.RemoveRoomMember(event.RoomID, client)

	if len(ids) == 0 {
		return nil
	}

	roomUsers := general.RoomUsersWSEvent{
		Users: ids,
	}

	bytes, _ := wsevent.EncodeSocketEventWrapper(&roomUsers)

	u.hubManager.Emit(event.RoomID, bytes, func(id string) bool { return true })

	return nil
}

func NewUserLeftHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &userLeftHandler{hubManager: hubManager}
}
