package repository

import (
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
)

type VoiceSessionMetadataRepository interface {
	CreateSessionMetadata(event event.MeetingStartedEvent) (*dto.VoiceSessionDTO, error)
	AppendAudioURL(id, audioURL string) error
}
