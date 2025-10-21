package request

import "time"

type RecordMeetingRequest struct {
	RoomID     string    `json:"roomId"`
	FinishTime time.Time `json:"finishTime"`
}
