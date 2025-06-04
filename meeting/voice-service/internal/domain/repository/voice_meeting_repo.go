package repository

import (
	"time"
	"voice-service/internal/domain/dto"
	"voice-service/internal/domain/event"
)

type VoiceMeetingRepository interface {
	FinishSession(endedAt time.Time, id, audioURL string) (*dto.VoiceSessionDTO, error)
	CreateSession(event event.MeetingStartedEvent) (*dto.VoiceSessionDTO, error)
}
