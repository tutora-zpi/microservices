package interfaces

import (
	"context"
	"recorder-service/internal/domain/broker"
	"recorder-service/internal/domain/event"
)

type Broker interface {
	Publish(ctx context.Context, e event.Event, dest broker.Destination) error
	PublishMultiple(ctx context.Context, ev event.Event, destinations ...broker.Destination) error
	Consume(ctx context.Context, exchange string) error
}
