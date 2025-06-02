package model

import (
	"time"
	"voice-service/internal/domain/event"

	"github.com/lib/pq"
)

type VoiceSession struct {
	// ID is meeting ID to easily identify the voice meeting.
	ID         string         `gorm:"primaryKey;type:uuid" json:"id"`
	StartedAt  time.Time      `gorm:"type:timestamp" json:"startedAt"`
	EndedAt    *time.Time     `gorm:"type:timestamp" json:"endedAt"`
	MemberIDs  pq.StringArray `gorm:"type:text[]" json:"memberIDs"`
	IsFinished bool           `gorm:"type:boolean;default:false" json:"isFinished"`

	// recorded audio from meeting url
	AudioURL *string `gorm:"type:text" json:"audioURL"`
}

func NewVoiceSession(event event.MeetingStartedEvent) Model {
	start, err := time.Parse(time.RFC3339, event.StartedTime)
	if err != nil {
		start = time.Now()
	}

	ids := make([]string, len(event.Members))

	for i, member := range event.Members {
		ids[i] = member.ID
	}

	return &VoiceSession{
		ID:        event.MeetingID,
		StartedAt: start,
		MemberIDs: ids,
	}
}
