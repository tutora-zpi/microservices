package eventhandler

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/event"
	"chat-service/internal/domain/repository"
	"context"
	"encoding/json"
	"fmt"
)

type reactHandler struct {
	repo repository.MessageRepository
}

// Handle implements interfaces.EventHandler.
func (s *reactHandler) Handle(ctx context.Context, body []byte) error {
	var event event.ReactOnMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	err := s.repo.React(ctx, event)

	return err
}

func NewReactHandler(repo repository.MessageRepository) interfaces.EventHandler {
	return &reactHandler{repo: repo}
}
