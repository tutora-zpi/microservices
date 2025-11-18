package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"recorder-service/internal/app/interfaces/handler"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/ws_event/general"
	"recorder-service/internal/infrastructure/webrtc/writer"
)

type roomUsersHandler struct {
	botService    service.BotService
	writerFactory writer.WriterFactory
}

// Handle implements handler.EventHandler.
func (r *roomUsersHandler) Handle(ctx context.Context, body []byte) error {
	var evt general.RoomUsersWSEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode event: %v", err)
		return fmt.Errorf("failed to decode event: %w", err)
	}

	err := r.botService.UpdateBotPeers(evt)

	return err
}

func NewRoomUsersHandler(
	botService service.BotService,
	writerFactory writer.WriterFactory,
) handler.EventHandler {
	return &roomUsersHandler{botService: botService, writerFactory: writerFactory}
}
