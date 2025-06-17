package interfaces

import (
	"fmt"
	"log"
	"sync"
)

type Dispatcher interface {
	Register(pattern string, handler UseCaseHandler)
	Dispatch(pattern string, payload []byte) error
	GetHandler(pattern string) (func([]byte) error, bool)
}

type dispatcherImlp struct {
	handlers map[string]UseCaseHandler // A map of event names to their corresponding use case handlers.
	mu       sync.RWMutex              // A read-write lock to ensure safe concurrent access to the handlers map.
}

// Register implements Dispatcher.
func (d *dispatcherImlp) Register(pattern string, handler UseCaseHandler) {
	if _, ok := d.handlers[pattern]; !ok {
		log.Printf("Assigning handler for %s\n", pattern)
		d.handlers[pattern] = handler
	} else {
		log.Printf("%s has already registered handler\n", pattern)
	}
}

// GetHandler implements Dispatcher.
func (d *dispatcherImlp) GetHandler(pattern string) (func([]byte) error, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handler, exists := d.handlers[pattern]
	return handler.Exec, exists
}

// HandleEvent implements Dispatcher.
func (d *dispatcherImlp) Dispatch(pattern string, payload []byte) error {
	if f, ok := d.handlers[pattern]; ok {
		return f.Exec(payload)
	}

	return fmt.Errorf("failed to find appropriate handler for %s", pattern)
}

func NewDispatcher() Dispatcher {
	return &dispatcherImlp{
		handlers: make(map[string]UseCaseHandler),
	}
}
