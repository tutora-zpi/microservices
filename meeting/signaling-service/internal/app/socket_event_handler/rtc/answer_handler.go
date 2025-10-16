package rtc

import (
	"context"
	"encoding/json"
	"fmt"
	"signaling-service/internal/app/interfaces"
	wsevent "signaling-service/internal/domain/ws_event"
	"signaling-service/internal/domain/ws_event/rtc"
)

type answerHandler struct {
	hubManager interfaces.HubManager
}

// Handle implements interfaces.EventHandler.
func (a *answerHandler) Handle(ctx context.Context, body []byte, client interfaces.Client) error {
	var event rtc.AnswerEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode %s payload", event.Name())
	}

	wrapper := wsevent.SocketEventWrapper{
		Payload: body,
		Name:    event.Name(),
	}

	a.hubManager.EmitToClient(event.To, [][]byte{wrapper.ToBytes()})

	return nil
}

func NewAnswerHandler(hubManager interfaces.HubManager) interfaces.EventHandler {
	return &answerHandler{hubManager: hubManager}
}
