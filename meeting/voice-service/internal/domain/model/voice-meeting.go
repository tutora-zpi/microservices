package model

import (
	"time"
	"voice-service/internal/domain/event"

	"github.com/lib/pq"
)

type VoiceMeeting struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StartedAt  time.Time      `gorm:"type:timestamp" json:"staredAt"`
	EndedAt    *time.Time     `gorm:"type:timestamp" json:"endedAt"`
	MemberIDs  pq.StringArray `gorm:"type:text[]" json:"memberIds"`
	IsFinished bool           `gorm:"type:boolean;default:false" json:"isFinished"`
}

func NewVoiceMeeting(event event.MeetingStartedEvent) Model {
	start, err := time.Parse(time.RFC3339, event.StartedTime)
	if err != nil {
		start = time.Now()
	}

	ids := make([]string, len(event.Members))

	for i, member := range event.Members {
		ids[i] = member.ID
	}

	return &VoiceMeeting{
		ID:         event.MeetingID,
		StartedAt:  start,
		MemberIDs:  ids,
		IsFinished: false,
	}
}
