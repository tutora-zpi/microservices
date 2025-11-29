package general

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type AutoDisconnect struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}

func (a *AutoDisconnect) IsValid() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (a *AutoDisconnect) Name() string {
	return reflect.TypeOf(*a).Name()
}
