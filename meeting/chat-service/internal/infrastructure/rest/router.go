package rest

import (
	"chat-service/internal/infrastructure/rest/v1/handlers"
	"net/http"

	_ "chat-service/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h handlers.Handlable) *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	api.Handle("/docs", http.RedirectHandler("/api/v1/docs/", http.StatusSeeOther))

	chat := api.PathPrefix("/chat").Subrouter()

	chat.Handle("/{id}", h.IsAuth(http.HandlerFunc(h.FindChat))).Methods(http.MethodGet)
	chat.Handle("/{id}", h.IsAuth(http.HandlerFunc(h.DeleteChat))).Methods(http.MethodDelete)
	chat.Handle("/general", h.IsAuth(handlers.Validate((http.HandlerFunc(h.CreateGeneralChat))))).Methods(http.MethodPost)
	chat.Handle("/{id}/messages", h.IsAuth(http.HandlerFunc(h.FetchMoreMessages))).Methods(http.MethodGet)

	return router
}
