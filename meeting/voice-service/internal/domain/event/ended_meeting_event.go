package event

import (
	"time"
	"voice-service/internal/domain/dto"
)

type MeetingEndedEvent struct {
	MeetingID string        `json:"meetingID"`
	Members   []dto.UserDTO `json:"members"`
	EndedTime time.Time     `json:"endedTime"`
}
