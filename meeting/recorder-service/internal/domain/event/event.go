package event

import (
	"encoding/json"
	"log"
)

type EventWrapper struct {
	Pattern string `json:"pattern"`
	Data    any    `json:"data"`
}

type Event any

func NewEventWrapper(pattern string, event Event) EventWrapper {
	return EventWrapper{
		Pattern: pattern,
		Data:    event,
	}
}

func (e *EventWrapper) ToJson() ([]byte, error) {
	return json.Marshal(e)
}

func (e *EventWrapper) FromJson(body []byte) *EventWrapper {
	var decoded EventWrapper
	err := json.Unmarshal(body, &decoded)
	if err != nil {
		log.Printf("Failed to decode message: %v", err)
		return nil
	}

	return &decoded
}

func (e *EventWrapper) DecodeBody(dest any) error {
	payload, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(payload, &dest); err != nil {
		return err
	}

	return nil
}
