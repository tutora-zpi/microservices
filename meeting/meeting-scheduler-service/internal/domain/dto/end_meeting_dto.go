package dto

import "github.com/go-playground/validator/v10"

type EndMeetingDTO struct {
	MeetingID string    `json:"meetingID" validate:"required,uuid4"`
	Members   []UserDTO `json:"members" validate:"required,min=2,dive"`
}

func (e *EndMeetingDTO) IsValid() error {
	validate := validator.New()
	return validate.Struct(e)
}
