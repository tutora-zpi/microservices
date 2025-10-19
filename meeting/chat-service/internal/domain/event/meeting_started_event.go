package event

import (
	"chat-service/internal/domain/dto"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
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

func (m *MeetingStartedEvent) IsValid() error {
	v := validator.New()

	return v.Struct(m)
}

func (m *MeetingStartedEvent) GetMemeberIDs() []string {
	var result []string = make([]string, len(m.Members))
	for i, member := range m.Members {
		result[i] = member.ID
	}

	return result
}
