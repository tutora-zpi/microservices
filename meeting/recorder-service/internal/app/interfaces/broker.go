package interfaces

import (
	"context"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/infrastructure/config"
)

type Broker interface {
	Close()

	Publish(event event.EventWrapper, opts config.PublishOptions) error
	Consume(ctx context.Context, options config.ConsumeOptions, dispacher Dispatcher) error
}
