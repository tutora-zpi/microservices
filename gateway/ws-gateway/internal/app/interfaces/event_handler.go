package interfaces

import (
	"context"
)

type EventHandler interface {
	Handle(ctx context.Context, body []byte, client Client) error
}
