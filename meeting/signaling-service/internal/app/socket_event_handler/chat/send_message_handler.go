package chat

import (
	"context"
	"signaling-service/internal/app/interfaces"
)

type sendMessageHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (s *sendMessageHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	panic("unimplemented")
}

func NewSendMessageHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &sendMessageHandler{hubManager: hubManager}
}
