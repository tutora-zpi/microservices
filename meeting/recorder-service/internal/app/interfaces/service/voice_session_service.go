package service

import (
	"context"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/dto/request"
)

type VoiceSessionService interface {
	GetAudio(ctx context.Context, req request.GetAudioRequest) (*dto.GetAudioDTO, error)
	GetSessions(ctx context.Context, req request.FetchSessions) ([]*dto.VoiceSessionMetadataDTO, error)
}
