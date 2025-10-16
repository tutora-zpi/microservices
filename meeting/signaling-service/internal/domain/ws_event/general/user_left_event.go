package general

import "github.com/go-playground/validator/v10"

type UserLeftEvent struct {
	RoomID string `json:"roomId" validate:"required,uuid4"`
	UserID string `json:"userId,omitempty"`
}

func (u *UserLeftEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserLeftEvent) Name() string {
	return "user-left"
}
