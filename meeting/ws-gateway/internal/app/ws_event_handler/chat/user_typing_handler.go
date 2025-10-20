package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/chat"
)

type userTypingHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *userTypingHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event chat.UserTypingWSEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	wrapper := wsevent.SocketEventWrapper{
		Name:    event.Name(),
		Payload: body,
	}

	u.hubManager.Emit(event.ChatID, wrapper.ToBytes(), func(id string) bool { return id != event.UserTyperID })

	return nil
}

func NewUserTypingHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &userTypingHandler{hubManager: hubManager}
}
