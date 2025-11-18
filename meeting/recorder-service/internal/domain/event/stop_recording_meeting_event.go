package event

import (
	"reflect"
	"time"
)

type StopRecordingMeetingEvent struct {
	RoomID   string    `json:"roomId"`
	StopTime time.Time `json:"stopTime"`
}

func (s *StopRecordingMeetingEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
