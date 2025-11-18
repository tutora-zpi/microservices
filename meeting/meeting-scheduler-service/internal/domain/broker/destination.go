package broker

import "meeting-scheduler-service/internal/domain/event"

type Destination struct {
	Queue      string
	RoutingKey string
	Exchange   string
}

func NewExchangeDestination(event event.Event, exchange string) Destination {
	return Destination{
		Queue:      "",
		Exchange:   exchange,
		RoutingKey: event.Name(),
	}
}
