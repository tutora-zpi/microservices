package service

import (
	"context"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/dto/request"
)

type VoiceSessionService interface {
	GetSessions(ctx context.Context, req request.FetchSessions) ([]*dto.VoiceSessionMetadataDTO, error)
}
