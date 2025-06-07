package dto

import "github.com/go-playground/validator/v10"

type UserDTO struct {
	ID        string  `json:"id" validate:"required,uuid4"`
	AvatarURL *string `json:"avatarURL,omitempty" validate:"omitempty,url"`
	FirstName string  `json:"firstName" validate:"required,min=2"`
	LastName  string  `json:"lastName" validate:"required,min=2"`
}

func (dto *UserDTO) IsValid() error {
	v := validator.New()

	return v.Struct(dto)
}
