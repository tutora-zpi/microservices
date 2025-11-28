package event

import (
	"reflect"
	"strings"
)

type RecordingsUploaded struct {
	ClassID string   `json:"classId"`
	RoomID  string   `json:"meetingId"`
	Merged  string   `json:"merged"`
	Voices  []string `json:"voices"`
}

func NewRecordingsUploaded(keys []string, classID, roomID string) *RecordingsUploaded {
	var r RecordingsUploaded
	for _, key := range keys {
		if strings.Contains(key, "merged") {
			r.Merged = key
		} else {
			r.Voices = append(r.Voices, key)
		}
	}

	r.ClassID = classID
	r.RoomID = roomID

	return &r
}

func (r *RecordingsUploaded) Name() string {
	return reflect.TypeOf(*r).Name()
}
