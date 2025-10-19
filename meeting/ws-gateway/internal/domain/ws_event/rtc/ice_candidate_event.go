package rtc

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type IceCandidateEvent struct {
	Candidate json.RawMessage `json:"candidate" validate:"required"`
	From      string          `json:"from" validate:"required,uuid4"`
	To        string          `json:"to" validate:"required,uuid4"`
}

func (i *IceCandidateEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(i)
}

func (i *IceCandidateEvent) Name() string {
	return reflect.TypeOf(*i).Name()
}
