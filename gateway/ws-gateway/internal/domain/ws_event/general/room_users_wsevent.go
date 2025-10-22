package general

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type RoomUsersWSEvent struct {
	Users []string `json:"users"`
}

func (r *RoomUsersWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *RoomUsersWSEvent) Name() string {
	return reflect.TypeOf(*r).Name()
}
