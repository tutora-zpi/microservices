package service

import (
	"context"
	"log"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/dto/request"
	"recorder-service/internal/domain/repository"
)

type voiceSessionServiceImpl struct {
	repo repository.VoiceSessionMetadataRepository
}

// GetSessions implements interfaces.VoiceSessionService.
func (v *voiceSessionServiceImpl) GetSessions(ctx context.Context, req request.FetchSessions) ([]*dto.VoiceSessionMetadataDTO, error) {
	log.Printf("Getting sessions %d for: %s", req.Limit, req.ClassID)
	result, err := v.repo.FetchSessionMetadata(ctx, req.ClassID, req.Limit, req.LastFetchedMeetingID)
	return result, err
}

func NewVoiceSessionService(repo repository.VoiceSessionMetadataRepository) service.VoiceSessionService {
	return &voiceSessionServiceImpl{repo: repo}
}
