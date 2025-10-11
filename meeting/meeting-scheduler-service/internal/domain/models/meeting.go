package models

import (
	"encoding/json"
	"meeting-scheduler-service/internal/domain/dto"
	"time"
)

type Meeting struct {
	ClassID    string   `json:"classId"`
	MeetingID  string   `json:"meetingId"`
	Timestamp  int64    `json:"timestamp"`
	Title      string   `json:"title"`
	MembersIDs []string `json:"membersIDs"`
}

func (m *Meeting) Json() []byte {
	if bytes, err := json.Marshal(m); err != nil {
		return []byte{}
	} else {
		return bytes
	}
}

func (m *Meeting) DTO() *dto.MeetingDTO {
	start := time.Unix(m.Timestamp, 0).UTC().Truncate(time.Minute)

	var members []dto.UserDTO = make([]dto.UserDTO, len(m.MembersIDs))
	for i, id := range m.MembersIDs {
		members[i] = dto.UserDTO{ID: id}
	}

	return &dto.MeetingDTO{
		MeetingID:   m.MeetingID,
		StartedDate: &start,
		Title:       m.Title,
		Members:     members,
	}
}
