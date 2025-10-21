package rest

import (
	"net/http"
	"recorder-service/internal/infrastructure/handlers"
	"recorder-service/internal/infrastructure/ws"

	"github.com/gorilla/mux"
)

func NewRouter(gateway ws.Gateway) *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	// swagger

	router.HandleFunc("/ws", gateway.Handle)

	api := router.PathPrefix("/api/v1/").Subrouter()

	api.HandleFunc("/health", handlers.HandleHealth)

	return router
}
