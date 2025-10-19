package repository

import (
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/event"
	"context"
)

type MessageRepository interface {
	Save(ctx context.Context, event event.SendMessageEvent) error
	Reply(ctx context.Context, event event.ReplyOnMessageEvent) error
	React(ctx context.Context, event event.ReactMessageOnEvent) error
	Delete(ctx context.Context, dto requests.DeleteMessage) error
	FindMore(ctx context.Context, dto requests.GetMoreMessages) ([]*dto.MessageDTO, error)
}
