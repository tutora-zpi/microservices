package rtc

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type OfferEvent struct {
	Offer json.RawMessage `json:"offer" validate:"required"`
	From  string          `json:"from" validate:"required,uuid4"`
	To    string          `json:"to" validate:"required,uuid4"`
}

func (o *OfferEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(o)
}

func (o *OfferEvent) Name() string {
	return reflect.TypeOf(*o).Name()
}
