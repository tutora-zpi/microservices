package service

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/repository"
	"context"
	"log"
)

type messageServiceImpl struct {
	repo repository.MessageRepository
}

// GetMoreMessages implements interfaces.MessageService.
func (m *messageServiceImpl) GetMoreMessages(ctx context.Context, req requests.GetMoreMessages) ([]*dto.MessageDTO, error) {
	log.Println("Getting more messages...")
	return m.repo.FindMore(ctx, req)
}

func NewMessageService(repo repository.MessageRepository) interfaces.MessageService {
	return &messageServiceImpl{repo: repo}
}
