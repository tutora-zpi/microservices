package event

import (
	"recorder-service/internal/domain/dto"
	"time"
)

type MeetingStartedEvent struct {
	ClassID     string        `json:"classId"`
	MeetingID   string        `json:"meetingId"`
	Members     []dto.UserDTO `json:"members"`
	StartedTime time.Time     `json:"startedTime"` // ISO 8601 format
	FinishTime  time.Time     `json:"finishTime"`
}
