package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"reflect"
	"time"
)

type MeetingEndedEvent struct {
	MeetingID string        `json:"meetingID"`
	Members   []dto.UserDTO `json:"members"`
	EndedTime time.Time     `json:"endedTime"`
}

func NewMeetingEndedEvent(dto dto.EndMeetingDTO) *EventWrapper {
	event := &MeetingEndedEvent{
		MeetingID: dto.MeetingID,
		Members:   dto.Members,
		EndedTime: time.Now(),
	}

	name := reflect.TypeOf(*event).Name()
	return NewEventWrapper(name, *event)
}
