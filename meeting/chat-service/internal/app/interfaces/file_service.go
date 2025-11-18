package interfaces

import (
	"chat-service/internal/domain/metadata"
	"context"
)

type FileService interface {
	Save(ctx context.Context, file *metadata.FileMetadata) (string, error)
}
