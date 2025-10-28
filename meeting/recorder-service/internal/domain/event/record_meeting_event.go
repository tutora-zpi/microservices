package event

import (
	"reflect"
	"time"
)

type RecordMeetingEvent struct {
	RoomID     string    `json:"roomId"`
	FinishTime time.Time `json:"finishTime"`
}

func (r *RecordMeetingEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
