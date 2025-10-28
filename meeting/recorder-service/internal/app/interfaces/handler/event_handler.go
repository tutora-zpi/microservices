package handler

import "context"

type EventHandler interface {
	Handle(ctx context.Context, body []byte) error
}
