package board

import "github.com/go-playground/validator/v10"

type BoardUpdateEvent struct {
	MeetingID string `json:"meetingId" validate:"reiqured,uuid4"`
	BoardSyncEvent
}

func (b *BoardUpdateEvent) IsValid() error {
	validate := validator.New()

	return validate.Struct(b)
}

func (b *BoardUpdateEvent) Type() string {
	return "board:update"
}

func (u *BoardUpdateEvent) Name() string {
	return u.Type()
}
