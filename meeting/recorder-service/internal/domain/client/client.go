package client

import (
	"context"
)

type Client interface {
	Send(msg []byte) error
	GetBotID() string
	Close() error
	IsConnected() bool
	Connect(ctx context.Context) error
}
