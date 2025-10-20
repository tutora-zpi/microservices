package eventhandler

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/event"
	"chat-service/internal/domain/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type meetingStartedHandler struct {
	repo repository.ChatRepository
}

// Handle implements interfaces.EventHandler.
func (s *meetingStartedHandler) Handle(ctx context.Context, body []byte) error {
	var event event.MeetingStartedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}
	log.Println(event.MeetingID)
	_, err := s.repo.Save(ctx, event.GetMemeberIDs(), event.MeetingID)
	return err
}

func NewMeetingStartedHandler(repo repository.ChatRepository) interfaces.EventHandler {
	return &meetingStartedHandler{repo: repo}
}
