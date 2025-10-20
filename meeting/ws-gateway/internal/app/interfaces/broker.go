package interfaces

import (
	"context"
	"ws-gateway/internal/domain/broker"
	"ws-gateway/internal/domain/event"
)

type Broker interface {
	Publish(ctx context.Context, e event.Event, dest broker.Destination) error
	PublishMultiple(ctx context.Context, ev event.Event, destinations ...broker.Destination) error
	Consume(ctx context.Context, exchange string) error
	Close()
}
