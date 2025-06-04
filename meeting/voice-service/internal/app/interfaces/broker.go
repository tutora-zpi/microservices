package interfaces

import "voice-service/internal/domain/event"

type Broker interface {
	Close()

	Publish(event event.EventWrapper) error
	Consume(event event.EventWrapper) error
}
