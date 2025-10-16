package interfaces

import (
	"context"
	"signaling-service/internal/domain/broker"
	"signaling-service/internal/domain/event"
)

type Broker interface {
	Publish(ctx context.Context, e event.Event, dest broker.Destination) error
	PublishMultiple(ctx context.Context, ev event.Event, destinations ...broker.Destination) error
	Close()
}
