package event

import (
	"reflect"
	"ws-gateway/internal/domain/dto"
)

type MeetingEndedEvent struct {
	MeetingID    string        `json:"meetingId"`
	EndTimestamp int64         `json:"endTimestamp"`
	Members      []dto.UserDTO `json:"members"`
}

func (m *MeetingEndedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}
