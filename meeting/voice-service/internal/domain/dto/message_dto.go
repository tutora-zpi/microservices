package dto

import (
	"github.com/go-playground/validator/v10"
)

// WsMessageDTO represents a WebSocket message.
// swagger:model WsMessageDTO
type WsMessageDTO struct {
	// Action performed in the message, e.g. "recording"
	// required: true
	// example: recording
	Action string `json:"action" validate:"required"`

	// Unique meeting identifier (UUID v4)
	// required: true
	// example: "123e4567-e89b-12d3-a456-426614174000"
	MeetingID string `json:"meeting_id" validate:"required,uuid4"`

	// Unique user identifier (UUID v4)
	// required: true
	// example: "987e6543-e21b-12d3-a456-426655440000"
	UserID string `json:"user_id" validate:"required,uuid4"`
}

// Command defines available commands for recording
// swagger:model Command
type Command string

const (
	// Start command to start recording
	Start Command = "start"
	// Stop command to stop recording
	Stop Command = "stop"
)

// RecordMessageDTO extends WsMessageDTO with a recording command
// swagger:model RecordMessageDTO
type RecordMessageDTO struct {
	WsMessageDTO

	// Command to control recording: "start" or "stop"
	// required: true
	// enum: start,stop
	// example: start
	Command Command `json:"command" validate:"required,oneof=start stop"`
}

func (r *RecordMessageDTO) IsValid() error {
	validate := validator.New()
	return validate.Struct(r)
}
