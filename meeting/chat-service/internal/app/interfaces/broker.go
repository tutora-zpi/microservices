package interfaces

import (
	"chat-service/internal/domain/broker"
	"chat-service/internal/domain/event"
	"context"
)

type Broker interface {
	Publish(ctx context.Context, e event.Event, dest broker.Destination) error
	PublishMultiple(ctx context.Context, ev event.Event, destinations ...broker.Destination) error
	Consume(ctx context.Context, exchange string) error
}
