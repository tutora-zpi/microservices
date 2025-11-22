package repository

import (
	"context"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
)

type VoiceSessionMetadataRepository interface {
	CreateSessionMetadata(ctx context.Context, event event.MeetingStartedEvent) (*dto.VoiceSessionMetadataDTO, error)
	AppendAudioName(ctx context.Context, meetingID, audioName string) (*dto.VoiceSessionMetadataDTO, error)
	FetchSessionMetadata(ctx context.Context, meetingID string, limit int64, lastFetchedID *string) ([]*dto.VoiceSessionMetadataDTO, error)
}
