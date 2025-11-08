package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"recorder-service/internal/app/interfaces/handler"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/ws_event/rtc"
)

type iceCandidateHandler struct {
	botService service.BotService
}

// Handle implements handler.EventHandler.
func (i *iceCandidateHandler) Handle(ctx context.Context, body []byte) error {
	var evt rtc.IceCandidateWSEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode ICE candidate event: %v", err)
		return fmt.Errorf("invalid event payload")
	}

	// if err := evt.IsValid(); err != nil {
	// 	return err
	// }

	bot, ok := i.botService.GetBot(evt.RoomID)
	if !ok {
		return fmt.Errorf("bot not found for room %s", evt.RoomID)
	}

	if bot.ID() != evt.To {
		return fmt.Errorf("ICE candidate not for this bot (to=%s, bot=%s)", evt.To, bot.ID())
	}

	peer, ok := bot.GetPeer(evt.From)
	if !ok {
		return fmt.Errorf("peer %s not found", evt.From)
	}

	if err := peer.AddICECandidate(evt.Candidate); err != nil {
		return fmt.Errorf("failed to create ice candidate: %v", err)
	}

	return nil
}

func NewIceCandidateHandler(botService service.BotService) handler.EventHandler {
	return &iceCandidateHandler{botService: botService}
}
