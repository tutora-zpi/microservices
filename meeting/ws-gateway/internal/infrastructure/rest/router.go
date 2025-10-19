package rest

import (
	"net/http"
	"ws-gateway/internal/infrastructure/rest/v1/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(handlers handlers.Handlable) *mux.Router {
	router := mux.NewRouter()

	ws := router.PathPrefix("/ws").Subrouter()

	// ws.Handle("", handlers.IsAuth(http.HandlerFunc(handlers.WebSocketHandler)))
	ws.Handle("", http.HandlerFunc(handlers.WebSocketHandler))
	return router
}
