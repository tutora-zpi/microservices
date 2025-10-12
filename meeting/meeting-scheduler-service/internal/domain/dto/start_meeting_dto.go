package dto

import (
	"fmt"
	"time"

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

	// Finish Date - date and time when meeting finshes, use toISOString to cast your date
	// required: true
	FinishDate time.Time `json:"finishDate" validate:"required" example:"2025-10-10T12:50:05+02:00"`
}

func (dto *StartMeetingDTO) ConvertTimeToUTC() {
	dto.FinishDate = dto.FinishDate.UTC().Truncate(time.Minute)
}

func (dto *StartMeetingDTO) IsValid() error {
	v := validator.New()

	if err := v.Struct(dto); err != nil {
		return err
	}

	if dto.FinishDate.Before(time.Now()) {
		return fmt.Errorf("finish date must be in the future")
	}

	dto.ConvertTimeToUTC()

	for _, member := range dto.Members {
		if err := member.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
