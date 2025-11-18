package interfaces

import (
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"context"
)

type ChatService interface {
	CreateChat(ctx context.Context, req requests.CreateGeneralChat) (*dto.ChatDTO, error)
	DeleteChat(ctx context.Context, req requests.DeleteChat) error
	FindChat(ctx context.Context, req requests.GetChat) (*dto.ChatDTO, error)
	UpdateChatMember(ctx context.Context, req requests.UpdateChatMembers) error
}
