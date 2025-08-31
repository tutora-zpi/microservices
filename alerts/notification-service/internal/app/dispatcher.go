package app

import (
	"log"
	"notification-serivce/internal/app/interfaces"
)

type Dispatcher struct {
	handlers map[string]interfaces.EventHandler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]interfaces.EventHandler),
	}
}

func (d *Dispatcher) Register(queueName string, handler interfaces.EventHandler) {
	d.handlers[queueName] = handler
}

func (d *Dispatcher) HandleEvent(queueName string, msg []byte) error {
	log.Printf("Handling event from '%s'\n", queueName)

	handler, exists := d.handlers[queueName]
	if !exists {
		log.Printf("No handler found for event type: %s", queueName)
		return nil
	}

	return handler.Handle(msg)
}
