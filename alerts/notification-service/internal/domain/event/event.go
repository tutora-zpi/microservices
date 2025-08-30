package event

import "encoding/json"

type event any

type EventWrapper struct {
	Pattern string `json:"pattern"`
	Data    event  `json:"data"`
}

func (e *EventWrapper) ToJson() ([]byte, error) {
	return json.Marshal(e)
}

func NewEventWrapper(pattern string, event event) *EventWrapper {
	return &EventWrapper{
		Pattern: pattern,
		Data:    event,
	}
}
