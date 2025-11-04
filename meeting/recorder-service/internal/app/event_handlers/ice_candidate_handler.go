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
		log.Printf("Failed to decode event: %v", err)
		return fmt.Errorf("failed to decode event: %s", evt.Name())
	}

	cmd, err := evt.ToCommand()
	if err != nil {
		return err
	}

	bot, ok := i.botService.GetBot(evt.RoomID)
	if !ok {
		return fmt.Errorf("bot not exists")
	}

	if bot.ID() != evt.To {
		return fmt.Errorf("message is not addressed to bot %s", bot.Name())
	}

	err = bot.Client().AddIceCandidate(cmd.Candidate)

	if err != nil {
		log.Printf("Failed to add ice candidate to bot %s", bot.Name())
		return fmt.Errorf("failed to add ice candidate")
	}

	log.Printf("Added ICE candidate for room %s", cmd.RoomID)

	return nil
}

func NewIceCandidateHandler(botService service.BotService) handler.EventHandler {
	return &iceCandidateHandler{botService: botService}
}
