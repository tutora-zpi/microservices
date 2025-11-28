package models

import (
	"encoding/json"
	"meeting-scheduler-service/internal/domain/dto"
	"time"
)

type Meeting struct {
	ClassID               string   `json:"classId"`
	MeetingID             string   `json:"meetingId"`
	ActualStartTimestamp  int64    `json:"actualStartTimestamp"`
	PredictedEndTimestamp int64    `json:"predictedEndTimestamp"`
	Title                 string   `json:"title"`
	MemberIDs             []string `json:"memberIds"`
}

func (m *Meeting) ToBytes() []byte {
	if bytes, err := json.Marshal(m); err != nil {
		return []byte{}
	} else {
		return bytes
	}
}

func (m *Meeting) DTO() *dto.MeetingDTO {
	start := time.Unix(m.ActualStartTimestamp, 0).UTC().Truncate(time.Minute)
	finish := time.Unix(m.PredictedEndTimestamp, 0).UTC()

	var members []dto.UserDTO = make([]dto.UserDTO, len(m.MemberIDs))
	for i, id := range m.MemberIDs {
		members[i] = dto.UserDTO{ID: id}
	}

	return &dto.MeetingDTO{
		MeetingID:   m.MeetingID,
		StartedDate: &start,
		FinishDate:  finish,
		Title:       m.Title,
		Members:     members,
	}
}
