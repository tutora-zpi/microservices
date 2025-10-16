package handlers

import (
	"net/http"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/infrastructure/bus"
	security "signaling-service/internal/infrastructure/security/repository"

	"github.com/gorilla/websocket"
)

type Handlable interface {
	WebSocketHandler(w http.ResponseWriter, r *http.Request)
	IsAuth(next http.Handler) http.Handler
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
