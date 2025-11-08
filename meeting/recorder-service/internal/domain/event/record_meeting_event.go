package event

import (
	"reflect"
	"time"
)

type RecordMeetingEvent struct {
	ExpectedUserIDs []string  `json:"userIds"`
	RoomID          string    `json:"roomId"`
	FinishTime      time.Time `json:"finishTime"`
}

func (r *RecordMeetingEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
