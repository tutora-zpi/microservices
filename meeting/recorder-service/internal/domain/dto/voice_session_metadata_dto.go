package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// VoiceSessionMetadataDTO represents a voice session data transfer object.
// swagger:model VoiceSessionMetadataDTO
type VoiceSessionMetadataDTO struct {
	// Unique identifier of the session
	// required: true
	ID string `json:"id" validate:"required"`

	// Unique identifier of the voice session (UUID v4) in our case it will be MeetingID.
	// required: true
	// example: "123e4567-e89b-12d3-a456-426614174000"
	MeetingID string `json:"meetingId" validate:"required,uuid4"`

	// ClassID is a identifier of class
	// required: true
	// example: "123e4567-e89b-12d3-a456-426614174000"
	ClassID string `json:"classId" validate:"required,uuid4"`

	// Duration of the session in seconds.
	// required: true
	// minimum: 1
	// example: 3600
	Duration int64 `json:"duration" validate:"required,min=1"`

	// ISO8601 formatted start time of the session.
	// required: true
	// example: "2025-06-08T12:34:56Z"
	StartedAt *time.Time `json:"startedAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`

	// ISO8601 formatted end time of the session (optional).
	// example: "2025-06-08T13:34:56Z"
	// nullable: true
	EndedAt *time.Time `json:"endedAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`

	// List of member IDs participating in the session.
	// required: true
	// min length: 2
	// example: ["uuid1", "uuid2"]
	MemberIDs []string `json:"memberIDs" validate:"required,min=2,dive,required"`

	// URL to the recorded audio file (optional).
	// nullable: true
	AudioName *string `json:"audioName,omitempty" validate:"omitempty,url"`
}

func (v *VoiceSessionMetadataDTO) IsFinished() bool {
	return v.EndedAt != nil
}

func (dto *VoiceSessionMetadataDTO) IsValid() error {
	v := validator.New()
	return v.Struct(dto)
}
