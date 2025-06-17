package dto

import "github.com/go-playground/validator/v10"

// EndMeetingDTO represents the request body to end a meeting.
// swagger:model EndMeetingDTO
type EndMeetingDTO struct {
	// Meeting unique identifier (UUIDv4)
	// required: true
	MeetingID string `json:"meetingID" validate:"required,uuid4"`
	// Members who participated in the meeting (minimum 2)
	// required: true
	Members []UserDTO `json:"members" validate:"required,min=2,dive"`
}

func (e *EndMeetingDTO) IsValid() error {
	validate := validator.New()
	return validate.Struct(e)
}
