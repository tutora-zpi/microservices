package bus

import (
	"fmt"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/query"
)

type queryBusImpl struct {
	registry *HandlerRegistry[interfaces.QueryHandler]
}

func NewQueryBus() interfaces.QueryBus {
	return &queryBusImpl{
		registry: NewHandlerRegistry[interfaces.QueryHandler](),
	}
}

func (q *queryBusImpl) Register(query query.Query, handler interfaces.QueryHandler) {
	q.registry.Register(query.Name(), handler)
}

func (q *queryBusImpl) HandleQuery(query query.Query) (any, error) {
	handlers := q.registry.GetHandlers(query.Name())
	if len(handlers) == 0 {
		return nil, fmt.Errorf("no handler registered for pattern %s", query.Name())
	}
	if len(handlers) > 1 {
		return nil, fmt.Errorf("multiple handlers registered for query %s", query.Name())
	}

	return handlers[0].Execute(query)
}

func (q *queryBusImpl) AvailablePatterns() []string {
	return q.registry.Patterns()
}
