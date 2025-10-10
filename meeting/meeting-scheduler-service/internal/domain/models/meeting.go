package models

import (
	"encoding/json"
	"meeting-scheduler-service/internal/domain/dto"
	"time"
)

type Meeting struct {
	ClassID   string `json:"classId"`
	MeetingID string `json:"meetingId"`
	Timestamp int64  `json:"timestamp"`
	Title     string `json:"title"`
}

func (m *Meeting) ToJSON() []byte {
	if bytes, err := json.Marshal(m); err != nil {
		return []byte{}
	} else {
		return bytes
	}
}

func (m *Meeting) ToDTO() *dto.MeetingDTO {
	start := time.Unix(m.Timestamp, 0).UTC().Truncate(time.Minute)

	return &dto.MeetingDTO{
		MeetingID:   m.MeetingID,
		StartedDate: &start,
		Title:       m.Title,
		Members:     []dto.UserDTO{},
	}
}
