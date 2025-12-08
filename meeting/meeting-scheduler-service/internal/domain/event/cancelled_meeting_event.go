package event

import (
	"meeting-scheduler-service/internal/domain/dto"
	"reflect"
	"time"
)

type CancelledMeetingEvent struct {
	Title       string        `json:"title"`
	Receivers   []dto.UserDTO `json:"members"`
	StartedDate time.Time     `json:"startedDate"`
}

func (c *CancelledMeetingEvent) Name() string {
	return reflect.TypeOf(*c).Name()
}
