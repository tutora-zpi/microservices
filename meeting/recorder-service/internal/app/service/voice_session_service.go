package service

import (
	"context"
	"log"
	"recorder-service/internal/app/interfaces/service"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/dto/request"
	"recorder-service/internal/domain/repository"
	"recorder-service/internal/infrastructure/s3"
)

type voiceSessionServiceImpl struct {
	repo      repository.VoiceSessionMetadataRepository
	s3Service s3.S3Service
}

// GetAudio implements service.VoiceSessionService.
func (v *voiceSessionServiceImpl) GetAudio(ctx context.Context, req request.GetAudioRequest) (*dto.GetAudioDTO, error) {
	url, err := v.s3Service.GetPresignURL(ctx, req.Key())
	if err != nil {
		return nil, err
	}
	result := &dto.GetAudioDTO{UrlToAudio: url, GetAudioRequest: req}
	return result, nil
}

// GetSessions implements interfaces.VoiceSessionService.
func (v *voiceSessionServiceImpl) GetSessions(ctx context.Context, req request.FetchSessions) ([]*dto.VoiceSessionMetadataDTO, error) {
	log.Printf("Getting sessions %d for: %s", req.Limit, req.MeetingID)
	result, err := v.repo.FetchSessionMetadata(ctx, req.MeetingID, req.Limit, req.LastFetchedID)
	return result, err
}

func NewVoiceSessionService(repo repository.VoiceSessionMetadataRepository, s3Service s3.S3Service) service.VoiceSessionService {
	return &voiceSessionServiceImpl{repo: repo, s3Service: s3Service}
}
