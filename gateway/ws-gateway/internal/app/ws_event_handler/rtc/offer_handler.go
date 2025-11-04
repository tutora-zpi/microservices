package rtc

import (
	"context"
	"encoding/json"
	"fmt"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/rtc"
)

type offerHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (o *offerHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event rtc.OfferWSEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	wrapper := wsevent.SocketEventWrapper{
		Payload: body,
		Name:    event.Name(),
	}

	o.hubManager.EmitToClient(event.To, [][]byte{wrapper.ToBytes()})

	return nil
}

func NewOfferHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &offerHandler{hubManager: hubManager}
}
