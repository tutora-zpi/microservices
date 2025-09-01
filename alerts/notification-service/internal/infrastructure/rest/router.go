package rest

import (
	"net/http"
	"notification-serivce/internal/infrastructure/middleware"
	"notification-serivce/internal/infrastructure/rest/v1/handlers"
	"notification-serivce/internal/infrastructure/sse"

	"github.com/gorilla/mux"
)

func NewRouter(manager *sse.SSEManager) *mux.Router {
	router := mux.NewRouter()
	requestHandler := handlers.NewRequestHandler(manager)

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	// register more handlers
	api := router.PathPrefix("/api/v1/notification").Subrouter()

	api.Handle("/stream", middleware.IsAuth(http.HandlerFunc(requestHandler.StreamNotifications)))

	return router
}
