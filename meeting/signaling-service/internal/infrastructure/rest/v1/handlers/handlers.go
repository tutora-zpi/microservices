package handlers

import (
	"net/http"
	"signaling-service/internal/app/interfaces"
	"signaling-service/internal/infrastructure/bus"
	"signaling-service/internal/infrastructure/cache"

	"github.com/gorilla/websocket"
)

type Handlable interface {
	WebSocketHandler(w http.ResponseWriter, r *http.Request)
	IsAuth(next http.Handler) http.Handler
}

type handlers struct {
	dispatcher   bus.Dispachable
	hub          interfaces.HubManager
	tokenService cache.TokenService

	upgrader websocket.Upgrader
}

func NewHandlers(dispatcher bus.Dispachable, hub interfaces.HubManager, tokenService cache.TokenService) Handlable {
	return &handlers{
		dispatcher:   dispatcher,
		hub:          hub,
		tokenService: tokenService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
