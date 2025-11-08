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
	"recorder-service/internal/infrastructure/webrtc/writer"
)

type recordMeetingHandler struct {
	botService    service.BotService
	writerFactory writer.WriterFactory
}

// Handle implements interfaces.EventHandler.
func (r *recordMeetingHandler) Handle(ctx context.Context, body []byte) error {
	log.Print("RECORDING MEETING HANDLER")
	var evt event.RecordMeetingEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode event: %v", err)
		return fmt.Errorf("failed to decode event: %s", evt.Name())
	}

	bot, err := r.botService.AddNewBot(ctx, evt)
	if err != nil {
		log.Printf("Failed to add bot: %v", err)
		return nil
	}

	err = r.tryJoinRoom(bot, evt)
	if err != nil {
		return err
	}

	return nil
}

func (r *recordMeetingHandler) tryJoinRoom(bot bot.Bot, evt event.RecordMeetingEvent) error {
	joinEvent := &general.UserJoinedWSEvent{
		UserID: bot.ID(),
		RoomID: evt.RoomID,
	}

	log.Printf("Bot %s attempting to join room %s", bot.Name(), evt.RoomID)

	msg, err := wsevent.EncodeSocketEventWrapper(joinEvent)
	if err != nil {
		return fmt.Errorf("failed to encode join event: %w", err)
	}

	if err := bot.Client().Send(msg); err != nil {
		log.Printf("Failed to send join message: %v", err)
		return fmt.Errorf("failed to send join message: %w", err)
	}

	log.Printf("Bot %s sent join message to room %s", bot.Name(), evt.RoomID)

	return nil
}

func NewRecorderMeetingHandler(botService service.BotService, writerFactory writer.WriterFactory) handler.EventHandler {
	return &recordMeetingHandler{botService: botService, writerFactory: writerFactory}
}
