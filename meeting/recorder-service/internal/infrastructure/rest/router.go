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

	api := router.PathPrefix("/api/v1/").Subrouter()

	sessions := api.PathPrefix("/sessions").Subrouter()
	sessions.Handle("/{meeting_id}", h.IsAuth(http.HandlerFunc(h.FetchSessions))).Methods(http.MethodGet)
	sessions.Handle("/{meeting_id}/audio/{name}", h.IsAuth(http.HandlerFunc(h.GetAudio))).Methods(http.MethodGet)

	return router
}
