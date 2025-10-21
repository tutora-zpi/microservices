package general

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type UserLeftWSEvent struct {
	RoomID string `json:"roomId" validate:"required,uuid4"`
	UserID string `json:"userId,omitempty"`
}

func (u *UserLeftWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserLeftWSEvent) Name() string {
	return reflect.TypeOf(*u).Name()
}
