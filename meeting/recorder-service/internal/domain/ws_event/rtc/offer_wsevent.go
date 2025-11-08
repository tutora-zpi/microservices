package rtc

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pion/webrtc/v3"
)

type OfferWSEvent struct {
	Offer  webrtc.SessionDescription `json:"offer" validate:"required"`
	From   string                    `json:"from" validate:"required,uuid4"`
	To     string                    `json:"to" validate:"required,uuid4"`
	RoomID string                    `json:"roomId" validate:"required,uuid4"`
}

func (o *OfferWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(o)
}

func (o *OfferWSEvent) Name() string {
	return reflect.TypeOf(*o).Name()
}
