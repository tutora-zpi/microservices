package broker

import "signaling-service/internal/domain/event"

type Destination struct {
	Queue      string
	RoutingKey string
	Exchange   string
}

// Used for sending on multiple exchange
func NewMultipleDestination(event event.Event, exchanges ...string) []Destination {
	result := make([]Destination, 0, len(exchanges))

	for _, exchange := range exchanges {
		result = append(result, Destination{
			RoutingKey: event.Name(),
			Queue:      "",
			Exchange:   exchange,
		})
	}

	return result
}

func NewExchangeDestination(event event.Event, exchange string) Destination {
	return Destination{
		Queue:      "",
		Exchange:   exchange,
		RoutingKey: event.Name(),
	}
}
