package interfaces

import (
	"context"
)

type Broker interface {
	Consume(ctx context.Context, exchange string) error
	Close()
}
