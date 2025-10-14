package general

import "github.com/go-playground/validator/v10"

type UserLeftEvent struct {
	RoomID string `json:"roomId" validate:"required,uuid4"`
}

func (u *UserLeftEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserLeftEvent) Type() string {
	return "user-left"
}

func (u *UserLeftEvent) Name() string {
	return u.Type()
}
