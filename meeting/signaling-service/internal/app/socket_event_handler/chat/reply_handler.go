package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/domain/ws_event/chat"
)

type replyHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (u *replyHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var dest chat.ReplyOnMessageEvent
	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("failed to decode %s payload", dest.Name())
	}

	// TODO

	return nil
}

func NewReplyHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &replyHandler{hubManager: hubManager}
}
