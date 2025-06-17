package dto

import "github.com/go-playground/validator/v10"

// VoiceSessionDTO represents a voice session data transfer object.
// swagger:model VoiceSessionDTO
type VoiceSessionDTO struct {
	// Unique identifier of the voice session (UUID v4) in our case it will be MeetingID.
	// required: true
	// example: "123e4567-e89b-12d3-a456-426614174000"
	ID string `json:"id" validate:"required,uuid4"`

	// Duration of the session in seconds.
	// required: true
	// minimum: 1
	// example: 3600
	Duration int64 `json:"duration" validate:"required,min=1"`

	// ISO8601 formatted start time of the session.
	// required: true
	// example: "2025-06-08T12:34:56Z"
	StartedAt string `json:"startedAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`

	// ISO8601 formatted end time of the session (optional).
	// example: "2025-06-08T13:34:56Z"
	// nullable: true
	EndedAt *string `json:"endedAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`

	// List of member IDs participating in the session.
	// required: true
	// min length: 1
	// example: ["uuid1", "uuid2"]
	MemberIDs []string `json:"memberIDs" validate:"required,min=1,dive,required"`

	// URL to the recorded audio file (optional).
	// nullable: true
	AudioURL *string `json:"audioURL,omitempty" validate:"omitempty,url"`
}

func (v *VoiceSessionDTO) IsFinished() bool {
	return v.EndedAt != nil
}

func (dto *VoiceSessionDTO) IsValid() error {
	v := validator.New()
	return v.Struct(dto)
}
