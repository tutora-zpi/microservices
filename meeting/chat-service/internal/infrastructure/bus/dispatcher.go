package bus

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/event"
	"context"
	"log"
)

type Dispachable interface {
	Register(evt event.Event, handler interfaces.EventHandler)
	HandleEvent(ctx context.Context, queueName string, msg []byte) error
	AvailablePatterns() []string
}

type Dispatcher struct {
	registry *HandlerRegistry[interfaces.EventHandler]
}

func NewDispatcher() Dispachable {
	return &Dispatcher{
		registry: NewHandlerRegistry[interfaces.EventHandler](),
	}
}

func (d *Dispatcher) Register(evt event.Event, handler interfaces.EventHandler) {
	log.Printf("Registering [%s]", evt.Name())
	d.registry.Register(evt.Name(), handler)
}

func (d *Dispatcher) HandleEvent(ctx context.Context, queueName string, msg []byte) error {
	log.Printf("Handling event from '%s'\n", queueName)

	handlers := d.registry.GetHandlers(queueName)
	if len(handlers) == 0 {
		log.Printf("No handler found for event type: %s", queueName)
		return nil
	}

	for _, h := range handlers {
		if err := h.Handle(ctx, msg); err != nil {
			log.Printf("Error handling event %s: %v", queueName, err)
		}
	}

	return nil
}

func (d *Dispatcher) AvailablePatterns() []string {
	return d.registry.Patterns()
}
