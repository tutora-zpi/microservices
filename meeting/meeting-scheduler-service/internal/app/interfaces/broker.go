package interfaces

import "meeting-scheduler-service/internal/domain/event"

type Broker interface {
	Publish(e event.EventWrapper) error
	Consume(e event.EventWrapper) error

	Close()
}
