package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"recorder-service/internal/app/interfaces/handler"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/bot"
	"recorder-service/internal/domain/event"
	wsevent "recorder-service/internal/domain/ws_event"
	"recorder-service/internal/domain/ws_event/general"
)

type recordMeetingHandler struct {
	botService service.BotService
}

// Handle implements interfaces.EventHandler.
func (r *recordMeetingHandler) Handle(ctx context.Context, body []byte) error {
	var evt event.RecordMeetingEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode event: %v", err)
		return fmt.Errorf("failed to decode event: %s", evt.Name())
	}

	bot, err := r.botService.AddNewBot(ctx, evt)
	if err != nil {
		return nil
	}

	backCtx := context.Background()

	err = bot.Client().Connect(backCtx)
	if err != nil {
		return err
	}

	err = r.tryJoinRoom(bot, evt)
	if err != nil {
		return err
	}

	go bot.Recorder().StartRecording(backCtx, evt.RoomID)

	return nil
}

func (r *recordMeetingHandler) tryJoinRoom(bot bot.Bot, evt event.RecordMeetingEvent) error {
	joinEvent := &general.UserJoinedWSEvent{
		UserID: bot.ID(),
		RoomID: evt.RoomID,
	}

	msg, err := wsevent.EncodeSocketEventWrapper(joinEvent)
	if err != nil {
		return err
	}

	err = bot.Client().Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func NewRecorderMeetingHandler(botService service.BotService) handler.EventHandler {
	return &recordMeetingHandler{botService: botService}
}
