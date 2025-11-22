package rest

import (
	"net/http"
	_ "recorder-service/docs"

	"recorder-service/internal/infrastructure/rest/v1/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h handlers.Handler) *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(h.NotFound)

	router.NotFoundHandler = http.HandlerFunc(h.NotFoundHandler)

	api := router.PathPrefix("/api/v1/").Subrouter()

	api.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	api.Handle("/docs", http.RedirectHandler("/api/v1/docs/", http.StatusSeeOther))

	sessions := api.PathPrefix("/sessions").Subrouter()
	sessions.Handle("/{meeting_id}", h.IsAuth(http.HandlerFunc(h.FetchSessions))).Methods(http.MethodGet)
	sessions.Handle("/{meeting_id}/audio/{name}", h.IsAuth(http.HandlerFunc(h.GetAudio))).Methods(http.MethodGet)

	return router
}
