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

type answerHandler struct {
	botService service.BotService
}

// Handle implements interfaces.EventHandler.
func (a *answerHandler) Handle(ctx context.Context, body []byte) error {
	var evt rtc.AnswerWSEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode event: %v", err)
		return fmt.Errorf("failed to decode event: %s", evt.Name())
	}

	desc, err := evt.ToCommand()
	if err != nil {
		log.Printf("Failed to map to desc command: %v", err)
		return fmt.Errorf("invalid description command")
	}

	bot, ok := a.botService.GetBot(evt.RoomID)
	if !ok {
		return fmt.Errorf("bot does not exists")
	}

	if bot.ID() != evt.To {
		return fmt.Errorf("message is not adressed for bot %s", bot.Name())
	}

	err = bot.Client().SetRemoteDescription(desc.SDP)
	if err != nil {
		log.Printf("Failed to set remote desc: %v", err)
		return fmt.Errorf("failed to set remote description")
	}

	log.Printf("Successfully answered from: %s by %s", evt.From, bot.Name())

	return nil
}

func NewAnswerHandler(botService service.BotService) handler.EventHandler {
	return &answerHandler{botService: botService}
}
