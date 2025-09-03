package interfaces

import (
	"notification-serivce/internal/domain/query"
)

type QueryBus interface {
	Register(query query.Query, handler QueryHandler)
	HandleQuery(query query.Query) (any, error)
	AvailablePatterns() []string
}
