package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/config"
	"ws-gateway/internal/domain/broker"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/chat"
	"ws-gateway/internal/infrastructure/bus"
)

type reactHandler struct {
	hubManager  interfaces.HubManager
	eventBuffer bus.EventBuffer
	exchange    string
}

// Handle implements interfaces.EventHandler.
func (r *reactHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event chat.ReactOnMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	wrapper := wsevent.SocketEventWrapper{
		Name:    event.Name(),
		Payload: body,
	}

	go r.hubManager.Emit(event.ChatID, wrapper.ToBytes(), func(id string) bool { return true })

	r.eventBuffer.Add(&event, broker.NewExchangeDestination(&event, r.exchange))

	return nil
}

func NewReactHandler(hubManager interfaces.HubManager, eventBuffer bus.EventBuffer) interfaces.EventHandler {
	ex := os.Getenv(config.CHAT_EXCHANGE)
	return &reactHandler{hubManager: hubManager, eventBuffer: eventBuffer, exchange: ex}
}
