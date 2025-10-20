package interfaces

import (
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/event"
	"context"
)

type MessageService interface {
	GetMoreMessages(ctx context.Context, req requests.GetMoreMessages) ([]*dto.MessageDTO, error)
	SaveFileMessage(ctx context.Context, event event.SendMessageEvent) (*dto.MessageDTO, error)
}
