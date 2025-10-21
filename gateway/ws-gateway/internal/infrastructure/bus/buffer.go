package bus

import (
	"context"
	"log"
	"time"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/domain/broker"
	"ws-gateway/internal/domain/event"
)

type EventBuffer interface {
	Add(event event.Event, dest broker.Destination)
	Flush(ctx context.Context, buffer []Package) error
	Work(ctx context.Context)
	Close()
}

type Package struct {
	dest  broker.Destination
	event event.Event
}

type eventBufferImpl struct {
	broker interfaces.Broker

	buffer chan Package
	closed chan struct{}
}

// Close implements EventBuffer.
func (e *eventBufferImpl) Close() {
	select {
	case <-e.closed:
		log.Println("Buffer already closed...")
	default:
		log.Println("Closing event buffer...")
		close(e.closed)
	}
}

// Work implements EventBuffer.
func (e *eventBufferImpl) Work(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var packages []Package

	for {
		select {
		case pack, ok := <-e.buffer:
			if !ok {
				log.Println("Buffer closed, exiting Work loop")
				return
			}
			packages = append(packages, pack)

		case <-ticker.C:
			if len(packages) > 0 {
				batch := append([]Package(nil), packages...)
				packages = packages[:0]

				if err := e.Flush(ctx, batch); err != nil {
					log.Println("Flush error:", err)
				}
			}
		case <-ctx.Done():
			log.Println("Event buffer stopped via context")
			return

		case <-e.closed:
			log.Println("Closing event buffer...")
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
		closed: make(chan struct{}),
	}
}
