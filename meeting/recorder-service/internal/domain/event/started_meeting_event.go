package event

import "recorder-service/internal/domain/dto"

type MeetingStartedEvent struct {
	MeetingID   string        `json:"meetingID"`
	Members     []dto.UserDTO `json:"members"`
	StartedTime string        `json:"startedTime"` // ISO 8601 format
}
