package models

import (
	"encoding/json"
	"meeting-scheduler-service/internal/domain/dto"
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
	return &dto.MeetingDTO{
		MeetingID: m.MeetingID,
		Timestamp: &m.Timestamp,
		Title:     m.Title,
		Members:   []dto.UserDTO{},
	}
}
