package event

import (
	"recorder-service/internal/domain/dto"
	"time"
)

type MeetingEndedEvent struct {
	MeetingID string        `json:"meetingID"`
	Members   []dto.UserDTO `json:"members"`
	EndedTime time.Time     `json:"endedTime"`
}
