package interfaces

import "notification-serivce/internal/domain/event"

type Dispatcher interface {
	Register(event event.Event, handler EventHandler)
	HandleEvent(pattern string, msg []byte) error
	AvailablePatterns() []string
}
