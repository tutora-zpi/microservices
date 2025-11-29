package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"reflect"
	"time"
)

type MeetingEndedEvent struct {
	MeetingID    string        `json:"meetingId"`
	ClassID      string        `json:"classId"`
	EndTimestamp int64         `json:"endTimestamp"`
	Members      []dto.UserDTO `json:"members"`
}

func NewMeetingEndedEvent(dto dto.EndMeetingDTO, members []dto.UserDTO) *MeetingEndedEvent {
	event := &MeetingEndedEvent{
		MeetingID:    dto.MeetingID,
		EndTimestamp: time.Now().UTC().Unix(),
		Members:      members,
		ClassID:      dto.ClassID,
	}

	return event
}

func (m *MeetingEndedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}
