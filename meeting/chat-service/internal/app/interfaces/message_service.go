package interfaces

import (
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"context"
)

type MessageService interface {
	GetMoreMessages(ctx context.Context, req requests.GetMoreMessages) ([]*dto.MessageDTO, error)
}
