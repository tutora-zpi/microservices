package wsevent

import (
	"encoding/json"
	"fmt"
	"signaling-service/internal/domain/event"
	"signaling-service/internal/domain/validator"
)

type SocketEvent interface {
	Type() string

	validator.Validable
	event.Event
}

type SocketEventWrapper struct {
	Type string `json:"type"`

	// socket events name
	Name string `json:"name"`

	Payload []byte `json:"data"`
}

func DecodeSocketEventWrapper(payload []byte) (*SocketEventWrapper, error) {
	var dest SocketEventWrapper

	if err := json.Unmarshal(payload, &dest); err != nil {
		return nil, fmt.Errorf("failed to decode")
	}

	return &dest, nil
}
