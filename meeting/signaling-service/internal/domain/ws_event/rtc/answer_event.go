package rtc

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type AnswerEvent struct {
	Answer json.RawMessage `json:"answer" validate:"required"`
	From   string          `json:"from" validate:"required,uuid4"`
	To     string          `json:"to" validate:"required,uuid4"`
}

func (a *AnswerEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (a *AnswerEvent) Name() string {
	return "answer"
}
