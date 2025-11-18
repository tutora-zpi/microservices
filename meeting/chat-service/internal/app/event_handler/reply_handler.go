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

type replyHandler struct {
	repo repository.MessageRepository
}

// Handle implements interfaces.EventHandler.
func (s *replyHandler) Handle(ctx context.Context, body []byte) error {
	var event event.ReplyOnMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}
	log.Printf("Handling [%s]", event.Name())

	err := s.repo.Reply(ctx, event)

	return err
}

func NewReplyHandler(repo repository.MessageRepository) interfaces.EventHandler {
	return &replyHandler{repo: repo}
}
