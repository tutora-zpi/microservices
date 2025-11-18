package client

import (
	"context"
)

type Client interface {
	Send(msg []byte) error
	SetBotID(botID string)
	Close() error
	IsConnected() bool
	Connect(ctx context.Context) error
}
