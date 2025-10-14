package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/domain/ws_event/chat"
)

type reactHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *reactHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var dest chat.ReactOnMessageEvent
	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("failed to decode %s payload", dest.Name())
	}

	// u.hubManager.Emit(dest.ChatID, websocket.TextMessage, body, func(id string) bool { return id != dest.UserTyperID })

	return nil
}

func NewReactHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &reactHandler{hubManager: hubManager}
}
