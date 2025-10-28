package rest

import (
	"net/http"

	"recorder-service/internal/infrastructure/rest/v1/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(h handlers.Handler) *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(h.NotFound)

	// swagger

	// router.HandleFunc("/ws", gateway.Handle)

	api := router.PathPrefix("/api/v1/").Subrouter()

	sessions := api.PathPrefix("/sessions").Subrouter()
	sessions.Handle("/{meeting_id}", h.IsAuth(http.HandlerFunc(h.FetchSessions))).Methods(http.MethodGet)

	return router
}
