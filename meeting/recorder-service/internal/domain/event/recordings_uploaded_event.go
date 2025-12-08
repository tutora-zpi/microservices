package event

import (
	"reflect"
	"strings"
)

type RecordingsUploaded struct {
	ClassID   string   `json:"classId"`
	RoomID    string   `json:"meetingId"`
	Merged    string   `json:"merged"`
	Voices    []string `json:"voices"`
	MemberIDs []string `json:"memberIds"`
}

func NewRecordingsUploaded(keys []string, classID, roomID string, memberIds []string) *RecordingsUploaded {
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
	r.MemberIDs = memberIds

	return &r
}

func (r *RecordingsUploaded) Name() string {
	return reflect.TypeOf(*r).Name()
}
