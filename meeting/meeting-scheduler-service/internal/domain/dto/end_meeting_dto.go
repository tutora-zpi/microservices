package dto

import "github.com/go-playground/validator/v10"

// EndMeetingDTO represents the request body to end a meeting.
// swagger:model EndMeetingDTO
type EndMeetingDTO struct {
	// Meeting unique identifier (UUIDv4)
	// required: true
	MeetingID string `json:"meetingID" validate:"required,uuid4"`
	// Class id - where meeting will be started (UUIDv4)
	// required: true
	ClassID string `json:"classId" validate:"required,uuid4"`
}

func (e *EndMeetingDTO) IsValid() error {
	validate := validator.New()
	return validate.Struct(e)
}
