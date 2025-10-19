package service

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/repository"
	"context"
	"log"
)

type chatServiceImpl struct {
	repo repository.ChatRepository
}

// CreateChat implements interfaces.ChatService.
func (c *chatServiceImpl) CreateChat(ctx context.Context, req requests.CreateGeneralChat) (*dto.ChatDTO, error) {
	log.Println("Creating chat...")
	return c.repo.Save(ctx, req.MemberIDs, req.ClassID)
}

// DeleteChat implements interfaces.ChatService.
func (c *chatServiceImpl) DeleteChat(ctx context.Context, req requests.DeleteChat) error {
	log.Println("Deleting chat...")
	return c.repo.Delete(ctx, req)
}

// FindChat implements interfaces.ChatService.
func (c *chatServiceImpl) FindChat(ctx context.Context, req requests.GetChat) (*dto.ChatDTO, error) {
	log.Println("Finding chat...")
	return c.repo.Find(ctx, req)
}

func NewChatService(repo repository.ChatRepository) interfaces.ChatService {
	return &chatServiceImpl{repo: repo}
}
