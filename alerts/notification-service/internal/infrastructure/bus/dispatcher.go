package bus

import (
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/event"
)

type Dispatcher struct {
	registry *HandlerRegistry[interfaces.EventHandler]
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		registry: NewHandlerRegistry[interfaces.EventHandler](),
	}
}

func (d *Dispatcher) Register(evt event.Event, handler interfaces.EventHandler) {
	d.registry.Register(evt.Name(), handler)
}

func (d *Dispatcher) HandleEvent(queueName string, msg []byte) error {
	log.Printf("Handling event from '%s'\n", queueName)

	handlers := d.registry.GetHandlers(queueName)
	if len(handlers) == 0 {
		log.Printf("No handler found for event type: %s", queueName)
		return nil
	}

	for _, h := range handlers {
		if err := h.Handle(msg); err != nil {
			log.Printf("Error handling event %s: %v", queueName, err)
		}
	}

	return nil
}

func (d *Dispatcher) AvailablePatterns() []string {
	return d.registry.Patterns()
}
