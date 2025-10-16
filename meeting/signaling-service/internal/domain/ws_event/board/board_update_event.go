package board

import "github.com/go-playground/validator/v10"

type BoardUpdateEvent struct {
	MeetingID string         `json:"meetingId" validate:"reiqured,uuid4"`
	Data      map[string]any `json:"data"`
}

func (b *BoardUpdateEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(b)
}

func (b *BoardUpdateEvent) Name() string {
	return "board:update"
}
