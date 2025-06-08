package dto

import "github.com/go-playground/validator/v10"

// UserDTO represents a user participating in a meeting.
// swagger:model UserDTO
type UserDTO struct {
	// User unique identifier (UUIDv4)
	// required: true
	ID string `json:"id" validate:"required,uuid4"`
	// URL to the user's avatar image
	// required: false
	AvatarURL *string `json:"avatarURL,omitempty" validate:"omitempty,url"`
	// User's first name
	// required: true
	FirstName string `json:"firstName" validate:"required,min=2"`
	// User's last name
	// required: true
	LastName string `json:"lastName" validate:"required,min=2"`
}

func (dto *UserDTO) IsValid() error {
	v := validator.New()

	return v.Struct(dto)
}
