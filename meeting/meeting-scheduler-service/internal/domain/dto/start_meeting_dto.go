package dto

import (
	"github.com/go-playground/validator/v10"
)

// StartMeetingDTO represents the request body to start a meeting.
// swagger:model StartMeetingDTO
type StartMeetingDTO struct {
	// Members participating in the meeting (minimum 2)
	// required: true
	Members []UserDTO `json:"members" validate:"required,min=2,dive"`

	// Class id - where meeting will be started (UUIDv4)
	// required: true
	ClassID string `json:"classId" validate:"required,uuid4"`

	// The title of the class eg. C++ Object oriented: pointers
	// reqiured: true
	Title string `json:"title" validate:"required"`
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
