package general

import (
	"context"
	"encoding/json"
	"fmt"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/domain/ws_event/general"

	"github.com/gorilla/websocket"
)

type userLeftHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *userLeftHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var dest general.UserLeftEvent

	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("failed to decode: %v", err)
	}

	u.hubManager.RemoveMeetingMemeber(dest.RoomID, client)

	welcomeMsg := fmt.Appendf(nil, "%s has left room: %s", client.ID(), dest.RoomID)

	u.hubManager.Emit(dest.RoomID, websocket.TextMessage, welcomeMsg, func(id string) bool { return true })

	client.GetConnection().Close()

	return nil
}

func NewUserLeftHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &userLeftHandler{hubManager: hubManager}
}
