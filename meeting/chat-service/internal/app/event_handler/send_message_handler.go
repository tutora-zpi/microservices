package eventhandler

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/event"
	"chat-service/internal/domain/repository"
	"context"
	"encoding/json"
	"fmt"
)

type sendMessageHandlerImpl struct {
	repo repository.MessageRepository
}

// Handle implements interfaces.EventHandler.
func (s *sendMessageHandlerImpl) Handle(ctx context.Context, body []byte) error {
	var event event.SendMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	_, err := s.repo.Save(ctx, event)

	return err
}

func NewSendMessageHandler(repo repository.MessageRepository) interfaces.EventHandler {
	return &sendMessageHandlerImpl{repo: repo}
}
