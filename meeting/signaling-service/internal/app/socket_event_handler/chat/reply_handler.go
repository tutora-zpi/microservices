package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/config"
	"signaling-service/internal/domain/broker"
	"signaling-service/internal/domain/ws_event/chat"
	"signaling-service/internal/infrastructure/bus"
)

type replyHandler struct {
	hubManager  interfaces.HubManager
	eventBuffer bus.EventBuffer
	exchange    string
}

// Handle implements interfaces.EventHandler.
func (r *replyHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event chat.ReplyOnMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	go r.hubManager.Emit(event.ChatID, body, func(id string) bool { return true })

	r.eventBuffer.Add(&event, broker.NewExchangeDestination(&event, r.exchange))

	return nil
}

func NewReplyHandler(hubManager interfaces.HubManager, eventBuffer bus.EventBuffer) interfaces.EventHandler {
	ex := os.Getenv(config.CHAT_EXCHANGE)
	return &replyHandler{hubManager: hubManager, eventBuffer: eventBuffer, exchange: ex}
}
