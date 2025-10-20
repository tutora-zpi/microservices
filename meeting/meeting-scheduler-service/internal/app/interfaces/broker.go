package interfaces

import (
	"context"
	"meeting-scheduler-service/internal/domain/broker"
	"meeting-scheduler-service/internal/domain/event"
)

type Broker interface {
	Publish(ctx context.Context, e event.Event, dest broker.Destination) error
	PublishMultiple(ctx context.Context, ev event.Event, destinations ...broker.Destination) error
	Close()
}
