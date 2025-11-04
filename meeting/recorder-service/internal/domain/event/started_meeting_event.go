package event

import (
	"recorder-service/internal/domain/dto"
	"reflect"
	"time"
)

type MeetingStartedEvent struct {
	ClassID     string        `json:"classId"`
	MeetingID   string        `json:"meetingId"`
	Members     []dto.UserDTO `json:"members"`
	StartedTime time.Time     `json:"startedTime"` // ISO 8601 format
	FinishTime  time.Time     `json:"finishTime"`
}

func (m *MeetingStartedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}
