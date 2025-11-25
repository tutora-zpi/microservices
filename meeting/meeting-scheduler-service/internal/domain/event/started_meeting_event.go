package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"meeting-scheduler-service/internal/domain/models"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type MeetingStartedEvent struct {
	ClassID     string        `json:"classId"`
	MeetingID   string        `json:"meetingId"`
	Members     []dto.UserDTO `json:"members"`
	StartedTime time.Time     `json:"startedTime"` // ISO 8601 format
	FinishTime  time.Time     `json:"finishTime"`
}

func NewMeetingStartedEvent(dto dto.StartMeetingDTO) *MeetingStartedEvent {
	return &MeetingStartedEvent{
		MeetingID:   uuid.New().String(),
		Members:     dto.Members,
		StartedTime: time.Now().UTC().Truncate(time.Minute),
		FinishTime:  dto.FinishDate.UTC(),
		ClassID:     dto.ClassID,
	}
}

func (m *MeetingStartedEvent) Name() string {
	return reflect.TypeOf(*m).Name()
}

func (m *MeetingStartedEvent) NewMeeting(dto dto.StartMeetingDTO) *models.Meeting {
	ids := make([]string, len(dto.Members))
	for i, user := range dto.Members {
		ids[i] = user.ID
	}

	return &models.Meeting{
		MeetingID:             m.MeetingID,
		ActualStartTimestamp:  m.StartedTime.Unix(),
		PredictedEndTimestamp: m.FinishTime.Unix(),
		ClassID:               m.ClassID,
		Title:                 dto.Title,
		MemberIDs:             ids,
	}
}
