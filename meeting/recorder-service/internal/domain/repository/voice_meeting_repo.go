package repository

import (
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/event"
)

type VoiceMeetingRepository interface {
	FinishSession(id string) (*dto.VoiceSessionDTO, error)
	CreateSession(event event.MeetingStartedEvent) (*dto.VoiceSessionDTO, error)
	AppendAudioURL(id, audioURL string) error
}
