package general

import (
	"context"
	"encoding/json"
	"fmt"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/domain/ws_event/general"

	"github.com/gorilla/websocket"
)

type userJoinedHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *userJoinedHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var dest general.UserJoinedEvent

	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("failed to decode: %v", err)
	}

	u.hubManager.AddMeetingMember(dest.RoomID, client)

	welcomeMsg := fmt.Appendf(nil, "%s has joined to room: %s", client.ID(), dest.RoomID)

	u.hubManager.Emit(dest.RoomID, websocket.TextMessage, welcomeMsg)
	return nil
}

func NewUserJoinedHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &userJoinedHandler{hubManager: hubManager}
}
