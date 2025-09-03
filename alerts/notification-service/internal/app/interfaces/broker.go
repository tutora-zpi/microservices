package interfaces

import (
	"notification-serivce/internal/domain/event"
)

type Broker interface {
	Publish(e event.EventWrapper) error
	Consume() error

	Close()
}
