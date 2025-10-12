package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"reflect"
	"time"
)

type MeetingEndedEvent struct {
	MeetingID    string `json:"meetingID"`
	EndTimestamp int64  `json:"endTimestamp"`
}

func NewMeetingEndedEvent(dto dto.EndMeetingDTO) *MeetingEndedEvent {
	event := &MeetingEndedEvent{
		MeetingID:    dto.MeetingID,
		EndTimestamp: time.Now().UTC().Unix(),
	}

	return event
}

func (m *MeetingEndedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}
