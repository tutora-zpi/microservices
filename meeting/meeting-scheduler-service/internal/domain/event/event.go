package event

import (
	"encoding/json"
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

func (e *EventWrapper) ToJson() ([]byte, error) {
	return json.Marshal(e)
}
