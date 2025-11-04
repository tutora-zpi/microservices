package rtc

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pion/webrtc/v3"
)

type IceCandidateWSEvent struct {
	Candidate json.RawMessage `json:"candidate" validate:"required"`
	RoomID    string          `json:"roomId"`
	From      string          `json:"from" validate:"required,uuid4"`
	To        string          `json:"to" validate:"required,uuid4"`
}

func (i *IceCandidateWSEvent) IsValid() error {
	validate := validator.New()
	return validate.Struct(i)
}

func (i *IceCandidateWSEvent) Name() string {
	return reflect.TypeOf(*i).Name()
}

type AddIceCandidateCommand struct {
	RoomID    string
	Candidate webrtc.ICECandidateInit
}

func (e IceCandidateWSEvent) ToCommand() (*AddIceCandidateCommand, error) {
	var candidate webrtc.ICECandidateInit
	if err := json.Unmarshal(e.Candidate, &candidate); err != nil {
		return nil, fmt.Errorf("invalid ICE candidate payload: %w", err)
	}

	return &AddIceCandidateCommand{
		RoomID:    e.RoomID,
		Candidate: candidate,
	}, nil
}
