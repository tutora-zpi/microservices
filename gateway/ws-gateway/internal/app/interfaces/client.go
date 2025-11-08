package interfaces

import (
	"context"
)

type Client interface {
	ID() string
	GetConnection() Connection
	Listen(
		ctx context.Context,
		handler func(ctx context.Context, eventType string, msg []byte, client Client) error,
	)
	Close()
}

type Connection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, payload []byte) error
	Close()
}
