package wsevent

import (
	"encoding/json"
	"ws-gateway/internal/domain/event"
)

type SocketEventWrapper struct {
	Name string `json:"name,omitempty"`

	Payload json.RawMessage `json:"data,omitempty"`
}

func (s *SocketEventWrapper) ToBytes() []byte {
	if bytes, err := json.Marshal(s); err != nil {
		return nil
	} else {
		return bytes
	}
}

func DecodeSocketEventWrapper(payload []byte) (*SocketEventWrapper, error) {
	var dest SocketEventWrapper

	if err := json.Unmarshal(payload, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func EncodeSocketEventWrapper(event event.Event, name string) ([]byte, error) {
	bytes, err := json.Marshal(event)

	if err != nil {
		return nil, err
	}

	toEncode := SocketEventWrapper{
		Name:    name,
		Payload: bytes,
	}

	return json.Marshal(&toEncode)
}
