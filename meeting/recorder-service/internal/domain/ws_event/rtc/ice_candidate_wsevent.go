package rtc

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pion/webrtc/v3"
)

type IceCandidateWSEvent struct {
	Candidate webrtc.ICECandidateInit `json:"candidate" validate:"required"`
	RoomID    string                  `json:"roomId"`
	From      string                  `json:"from" validate:"required,uuid4"`
	To        string                  `json:"to" validate:"required,uuid4"`
}

func (i *IceCandidateWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(i)
}

func (i *IceCandidateWSEvent) Name() string {
	return reflect.TypeOf(*i).Name()
}
