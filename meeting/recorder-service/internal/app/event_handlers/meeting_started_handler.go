package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"recorder-service/internal/app/interfaces/handler"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/repository"
)

type meetingStartedHandler struct {
	repo repository.VoiceSessionMetadataRepository
}

// Handle implements handler.EventHandler.
func (m *meetingStartedHandler) Handle(ctx context.Context, body []byte) error {
	var evt event.MeetingStartedEvent

	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Failed to decode event: %v", err)
		return fmt.Errorf("failed to decode event: %s", evt.Name())
	}

	_, err := m.repo.CreateSessionMetadata(ctx, evt)

	log.Println("Session sucessfully created")
	return err
}

func NewMeetingStartedHandler(repo repository.VoiceSessionMetadataRepository) handler.EventHandler {
	return &meetingStartedHandler{repo: repo}
}
