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
	// socket events name
	Type string `json:"type"`

	Payload []byte `json:"data"`

	Metadata map[string]any `json:"meta,omitempty"`
}

func DecodeSocketEventWrapper(payload []byte) (*SocketEventWrapper, error) {
	var dest SocketEventWrapper

	if err := json.Unmarshal(payload, &dest); err != nil {
		return nil, fmt.Errorf("failed to decode")
	}

	return &dest, nil
}
