package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
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
		StartedTime: time.Now().UTC(),
	}

	return event
}

func (m *MeetingStartedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}

func (m *MeetingStartedEvent) NewMeeting(classID, title string) *models.Meeting {
	return &models.Meeting{
		MeetingID: m.MeetingID,
		Timestamp: m.StartedTime.Unix(),
		Title:     title,
		ClassID:   classID,
	}
}
