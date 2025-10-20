package rtc

import (
	"context"
	"encoding/json"
	"fmt"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"
	"ws-gateway/internal/domain/ws_event/rtc"
)

type iceCandidateHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (i *iceCandidateHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event rtc.IceCandidateWSEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	wrapper := wsevent.SocketEventWrapper{
		Payload: body,
		Name:    event.Name(),
	}

	i.hubManager.EmitToClient(event.To, [][]byte{wrapper.ToBytes()})

	return nil
}

func NewIceCandidateHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &iceCandidateHandler{hubManager: hubManager}
}
