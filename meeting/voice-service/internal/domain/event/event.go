package event

type EventWrapper struct {
	Pattern string `json:"pattern"`
	Data    Event  `json:"data"`
}

type Event interface {
}

func NewEventWrapper(pattern string, event Event) EventWrapper {
	return EventWrapper{
		Pattern: pattern,
		Data:    event,
	}
}
