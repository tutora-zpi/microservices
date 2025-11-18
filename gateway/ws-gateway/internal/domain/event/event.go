package event

import (
	"encoding/json"
	"log"
)

type EventWrapper struct {
	Pattern string `json:"pattern"`
	Data    Event  `json:"data"`
}

type Event interface {
	Name() string
}

func NewEventWrapper(pattern string, event Event) *EventWrapper {
	return &EventWrapper{
		Pattern: pattern,
		Data:    event,
	}
}

func (EventWrapper) DecodedEventWrapper(data []byte) (pattern string, decodedData []byte, err error) {
	type body struct {
		Pattern string          `json:"pattern"`
		Data    json.RawMessage `json:"data"`
	}

	var decodedBody body

	if err = json.Unmarshal(data, &decodedBody); err != nil {
		log.Println("Failed to decode data:", err)
		return "", nil, err
	}

	pattern = decodedBody.Pattern
	decodedData = decodedBody.Data

	return pattern, decodedData, nil
}

func (e *EventWrapper) ToBytes() []byte {
	if bytes, err := json.Marshal(e); err != nil {
		return []byte{}
	} else {
		return bytes
	}
}
