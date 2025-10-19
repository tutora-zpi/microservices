package general

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type RoomUsersEvent struct {
	Users []string `json:"users"`
}

func (r *RoomUsersEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *RoomUsersEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
