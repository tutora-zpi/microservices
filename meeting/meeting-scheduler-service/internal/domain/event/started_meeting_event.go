package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type MeetingStartedEvent struct {
	MeetingID   string        `json:"meetingId"`
	Members     []dto.UserDTO `json:"members"`
	StartedTime time.Time     `json:"startedTime"` // ISO 8601 format
}

func NewMeetingStartedEvent(dto dto.StartMeetingDTO) *MeetingStartedEvent {
	event := &MeetingStartedEvent{
		MeetingID:   uuid.New().String(),
		Members:     dto.Members,
		StartedTime: time.Now(),
	}

	return event
}

func (m *MeetingStartedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}
