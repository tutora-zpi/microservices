package bus

import (
	"context"
	"log"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/domain/broker"
	"signaling-service/internal/domain/event"
	"time"
)

type EventBuffer interface {
	Add(event event.Event, dest broker.Destination)
	Flush(ctx context.Context, buffer []Package) error
	Work(ctx context.Context)
}

type Package struct {
	dest  broker.Destination
	event event.Event
}

type eventBufferImpl struct {
	broker interfaces.Broker

	buffer chan Package
}

// Work implements EventBuffer.
func (e *eventBufferImpl) Work(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	packages := []Package{}

	for {
		select {
		case pack := <-e.buffer:
			log.Printf("Buffering new package: %v", pack)
			packages = append(packages, pack)
		case <-ticker.C:
			if len(packages) > 0 {
				if err := e.Flush(ctx, packages); err != nil {
					log.Println("Flush error:", err)
				}
				packages = packages[:0]
			}
		case <-ctx.Done():
			log.Println("Event buffer stopped")
			return
		}
	}
}

// Flush implements EventBuffer.
func (e *eventBufferImpl) Flush(ctx context.Context, buffer []Package) error {
	var err error
	for _, elem := range buffer {
		err = e.broker.Publish(ctx, elem.event, elem.dest)
	}

	return err
}

// Add implements EventBuffer.
func (e *eventBufferImpl) Add(event event.Event, dest broker.Destination) {
	newPackage := Package{
		dest:  dest,
		event: event,
	}

	e.buffer <- newPackage
}

func NewEventBuffer(broker interfaces.Broker) EventBuffer {
	return &eventBufferImpl{
		broker: broker,
		buffer: make(chan Package, 1000),
	}
}
