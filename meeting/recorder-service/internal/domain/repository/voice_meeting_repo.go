package repository

import (
	"context"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
)

type VoiceSessionMetadataRepository interface {
	CreateSessionMetadata(ctx context.Context, event event.MeetingStartedEvent) (*dto.VoiceSessionMetadataDTO, error)
	AppendAudioName(ctx context.Context, id, audioName string) error
	FetchSessionMetadata(ctx context.Context, classID string, limit int64, lastFetchedMeetingID *string) ([]*dto.VoiceSessionMetadataDTO, error)
}
