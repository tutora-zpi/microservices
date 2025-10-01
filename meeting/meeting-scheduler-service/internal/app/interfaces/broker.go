package interfaces

import "meeting-scheduler-service/internal/domain/event"

type Broker interface {
	Publish(e event.Event, exchangeChannels ...string) error

	Close()
}
