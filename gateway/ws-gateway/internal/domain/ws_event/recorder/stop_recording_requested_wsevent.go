package recorder

import (
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

type StopRecordingRequestedWSEvent struct {
	RoomID   string    `json:"roomId"`
	StopTime time.Time `json:"stopTime"`
	BotID    string    `json:"botId"`
}

func (s *StopRecordingRequestedWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (s *StopRecordingRequestedWSEvent) Name() string {
	return reflect.TypeOf(*s).Name()
}
