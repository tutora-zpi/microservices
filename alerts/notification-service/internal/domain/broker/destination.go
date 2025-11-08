package broker

import "notification-serivce/internal/domain/event"

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
