package bus

import (
	"context"
	"log"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/domain/event"
)

type Dispachable interface {
	Register(evt event.Event, handler interfaces.EventHandler)
	HandleEvent(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error
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
	log.Printf("Registering %s", evt.Name())
	d.registry.Register(evt.Name(), handler)
}

func (d *Dispatcher) HandleEvent(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error {
	log.Printf("Handling event from '%s'", eventType)

	handlers := d.registry.GetHandlers(eventType)
	if len(handlers) == 0 {
		log.Printf("No handler found for event type: %s", eventType)
		return nil
	}

	for _, h := range handlers {
		if err := h.Handle(ctx, msg, client); err != nil {
			log.Println(string(msg))
			log.Printf("Error handling event %s: %v", eventType, err)
		}
	}

	return nil
}

func (d *Dispatcher) AvailablePatterns() []string {
	return d.registry.Patterns()
}
