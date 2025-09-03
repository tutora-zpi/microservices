package sse

import (
	"context"
	"net/http"
	"notification-serivce/internal/app/interfaces"
	"time"
)

type ConnectionConfig struct {
	ClientID          string
	Writer            http.ResponseWriter
	Flusher           http.Flusher
	Channel           chan []byte
	Context           context.Context
	HeartbeatInterval time.Duration
	Manager           interfaces.NotificationManager
}

type SSEConfig struct {
	HeartbeatInterval time.Duration
	FrontendURL       string
	MaxConnections    int
	BufferEnabled     bool
	BufferSize        int
	BufferTTL         time.Duration
}
