package rtc

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pion/webrtc/v3"
)

type AnswerWSEvent struct {
	Answer webrtc.SessionDescription `json:"answer" validate:"required"`
	RoomID string                    `json:"roomId"`
	From   string                    `json:"from" validate:"required,uuid4"`
	To     string                    `json:"to" validate:"required,uuid4"`
}

func (a *AnswerWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (a *AnswerWSEvent) Name() string {
	return reflect.TypeOf(*a).Name()
}

type SetRemoteDescriptionCommand struct {
	RoomID string
	SDP    webrtc.SessionDescription
}

func (e AnswerWSEvent) ToCommand() (*SetRemoteDescriptionCommand, error) {
	var sdp webrtc.SessionDescription

	bytes, err := json.Marshal(e.Answer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal answer: %w", err)
	}

	if err := json.Unmarshal(bytes, &sdp); err != nil {
		return nil, fmt.Errorf("invalid SDP payload: %w", err)
	}

	return &SetRemoteDescriptionCommand{
		RoomID: e.RoomID,
		SDP:    sdp,
	}, nil
}
