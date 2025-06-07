package dto

import (
	"github.com/go-playground/validator/v10"
)

// StartMeetingDTO represents the request body to start a meeting.
// swagger:model StartMeetingDTO
type StartMeetingDTO struct {
	// Members participating in the meeting (minimum 2)
	// required: true
	Members []UserDTO `json:"members" validate:"required,min=2,dive,required"`
}

func (dto *StartMeetingDTO) IsValid() error {
	v := validator.New()

	if err := v.Struct(dto); err != nil {
		return err
	}

	for _, member := range dto.Members {
		if err := member.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
