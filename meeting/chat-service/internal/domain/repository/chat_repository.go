package repository

import (
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"context"
)

type ChatRepository interface {
	Find(ctx context.Context, dto requests.GetChat) (*dto.ChatDTO, error)
	Save(ctx context.Context, memberIDs []string, chatID string) (*dto.ChatDTO, error)
	Delete(ctx context.Context, event requests.DeleteChat) error
}
