package general

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type UserJoinedWSEvent struct {
	RoomID string `json:"roomId" validate:"required,uuid4"`
	UserID string `json:"userId,omitempty"`
}

func (u *UserJoinedWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserJoinedWSEvent) Name() string {
	return reflect.TypeOf(*u).Name()
}
