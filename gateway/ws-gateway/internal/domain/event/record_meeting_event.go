package event

import (
	"reflect"
	"time"
	"ws-gateway/internal/domain/ws_event/recorder"
)

type RecordMeetingEvent struct {
	ExpectedUserIDs []string  `json:"userIds"`
	RoomID          string    `json:"roomId"`
	FinishTime      time.Time `json:"finishTime"`
}

func (r *RecordMeetingEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}

func NewRecordMeetingEvent(e recorder.RecordRequestedWSEvent, userIDs []string) *RecordMeetingEvent {
	return &RecordMeetingEvent{
		ExpectedUserIDs: userIDs,
		RoomID:          e.RoomID,
		FinishTime:      e.FinishTime,
	}
}
