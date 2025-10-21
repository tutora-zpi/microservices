package board

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type BoardUpdateWSEvent struct {
	MeetingID string         `json:"meetingId" validate:"reiqured,uuid4"`
	Data      map[string]any `json:"data"`
}

func (b *BoardUpdateWSEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(b)
}

func (b *BoardUpdateWSEvent) Name() string {
	return reflect.TypeOf(*b).Name()
}
