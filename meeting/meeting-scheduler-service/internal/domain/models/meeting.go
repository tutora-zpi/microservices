package models

import (
	"encoding/json"
	"meeting-scheduler-service/internal/domain/dto"
	"time"
)

type Meeting struct {
	ClassID   string   `json:"classId"`
	MeetingID string   `json:"meetingId"`
	Timestamp int64    `json:"timestamp"`
	Title     string   `json:"title"`
	MemberIDs []string `json:"memberIds"`
}

func (m *Meeting) ToBytes() []byte {
	if bytes, err := json.Marshal(m); err != nil {
		return []byte{}
	} else {
		return bytes
	}
}

func (m *Meeting) DTO() *dto.MeetingDTO {
	start := time.Unix(m.Timestamp, 0).UTC().Truncate(time.Minute)

	var members []dto.UserDTO = make([]dto.UserDTO, len(m.MemberIDs))
	for i, id := range m.MemberIDs {
		members[i] = dto.UserDTO{ID: id}
	}

	return &dto.MeetingDTO{
		MeetingID:   m.MeetingID,
		StartedDate: &start,
		Title:       m.Title,
		Members:     members,
	}
}
