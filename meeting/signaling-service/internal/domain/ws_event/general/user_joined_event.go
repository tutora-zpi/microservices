package general

import "github.com/go-playground/validator/v10"

type UserJoinedEvent struct {
	RoomID string `json:"roomId" validate:"required,uuid4"`
	Token  string `json:"token"`
}

func (u *UserJoinedEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserJoinedEvent) Type() string {
	return "user-joined"
}

func (u *UserJoinedEvent) Name() string {
	return u.Type()
}
