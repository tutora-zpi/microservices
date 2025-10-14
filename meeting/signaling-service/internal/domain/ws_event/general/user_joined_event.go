package general

import "github.com/go-playground/validator/v10"

type UserJoinedEvent struct {
	RoomID string `json:"roomId" validate:"required,uuid4"`
}

func (u *UserJoinedEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserJoinedEvent) Name() string {
	return "user-joined"
}
