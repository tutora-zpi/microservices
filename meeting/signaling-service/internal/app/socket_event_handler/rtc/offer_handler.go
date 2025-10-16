package rtc

import (
	"context"
	"encoding/json"
	"fmt"
	"signaling-service/internal/app/interfaces"
	wsevent "signaling-service/internal/domain/ws_event"
	"signaling-service/internal/domain/ws_event/rtc"
)

type offerHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (o *offerHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event rtc.IceCandidateEvent
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
