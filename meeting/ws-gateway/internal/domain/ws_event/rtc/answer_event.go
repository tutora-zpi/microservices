package rtc

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type AnswerWSEvent struct {
	Answer json.RawMessage `json:"answer" validate:"required"`
	From   string          `json:"from" validate:"required,uuid4"`
	To     string          `json:"to" validate:"required,uuid4"`
}

func (a *AnswerWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (a *AnswerWSEvent) Name() string {
	return reflect.TypeOf(*a).Name()
}
