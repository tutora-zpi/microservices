package event

import "reflect"

type RecordingsUploaded struct {
	RoomID        string   `json:"roomId"`
	RecordingKeys []string `json:"recordingKeys"`
}

func (r *RecordingsUploaded) Name() string {
	return reflect.TypeOf(*r).Name()
}
