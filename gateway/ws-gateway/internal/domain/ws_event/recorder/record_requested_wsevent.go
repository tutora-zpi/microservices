package recorder

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

type RecordRequestedWSEvent struct {
	WhoRequestedID string    `json:"recordingRequesterId" validate:"required,uuid4"`
	RoomID         string    `json:"roomId"`
	FinishTime     time.Time `json:"finishTime"`
}

func (r *RecordRequestedWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *RecordRequestedWSEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}

func (r *RecordRequestedWSEvent) ToBytes() []byte {
	if bytes, err := json.Marshal(&r); err != nil {
		return nil
	} else {
		return bytes
	}
}

func NewRecordRequestedWSEventFromBytes(data []byte) *RecordRequestedWSEvent {
	var e RecordRequestedWSEvent

	if err := json.Unmarshal(data, &e); err != nil {
		return nil
	}

	return &e
}
