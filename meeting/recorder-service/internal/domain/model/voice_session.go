package model

import (
	"recorder-service/internal/domain/event"
	"time"

	"github.com/lib/pq"
)

type VoiceSessionMetadata struct {
	// ID is meeting ID to easily identify the voice meeting.
	MeetingID string `gorm:"primaryKey;not null" json:"meetingId"`
	ClassID   string `gorm:"not null;unique" json:"classId"`

	StartedAt time.Time      `gorm:"type:date;not null" json:"startedAt"`
	EndedAt   *time.Time     `gorm:"type:date;not null" json:"endedAt"`
	MemberIDs pq.StringArray `gorm:"text[];not null" json:"memberIds"`

	MergedAudioName *string `gorm:"not null" json:"audioName"`
}

func NewVoiceSession(event event.MeetingStartedEvent) *VoiceSessionMetadata {

	ids := make([]string, len(event.Members))

	for i, member := range event.Members {
		ids[i] = member.ID
	}

	return &VoiceSessionMetadata{
		MeetingID: event.MeetingID,
		ClassID:   event.ClassID,
		StartedAt: event.StartedTime,
		EndedAt:   &event.FinishTime,
		MemberIDs: ids,
	}
}
