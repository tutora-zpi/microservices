package handlers

import (
	"context"
	"net/http"
	"ws-gateway/internal/app/interfaces"
	"ws-gateway/internal/infrastructure/bus"
	security "ws-gateway/internal/infrastructure/security/repository"

	"github.com/gorilla/websocket"
)

type Handlable interface {
	WebSocketHandler(w http.ResponseWriter, r *http.Request)
	IsAuth(next http.Handler) http.Handler
	WithAuth(
		handler func(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error,
	) func(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error
}

type handlers struct {
	dispatcher bus.Dispachable
	hub        interfaces.HubManager
	tokenRepo  security.TokenRepository

	upgrader websocket.Upgrader
}

func NewHandlers(dispatcher bus.Dispachable, hub interfaces.HubManager, tokenService security.TokenRepository) Handlable {
	return &handlers{
		dispatcher: dispatcher,
		hub:        hub,
		tokenRepo:  tokenService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
