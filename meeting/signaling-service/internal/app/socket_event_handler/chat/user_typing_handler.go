package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/domain/ws_event/chat"
)

type userTypingHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *userTypingHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var dest chat.UserTypingEvent
	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("failed to decode %s payload", dest.Name())
	}

	u.hubManager.Emit(dest.ChatID, body, func(id string) bool { return id != dest.UserTyperID })

	return nil
}

func NewUserTypingHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &userTypingHandler{hubManager: hubManager}
}
