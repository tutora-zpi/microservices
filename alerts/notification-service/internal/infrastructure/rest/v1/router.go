package v1

import (
	"net/http"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/infrastructure/middleware"
	"notification-serivce/internal/infrastructure/rest/v1/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(manager interfaces.NotificationManager, queryBus interfaces.QueryBus) *mux.Router {
	router := mux.NewRouter()
	sseHandler := handlers.NewSSEHandler(manager, nil)
	httpHandler := handlers.NewHTTPHandler(queryBus)

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	api := router.PathPrefix("/api/v1/notification").Subrouter()

	api.Handle("/stream", middleware.IsAuth(http.HandlerFunc(sseHandler.StreamNotifications))).Methods(http.MethodGet)

	api.Handle("", middleware.IsAuth(http.HandlerFunc(httpHandler.FetchNotifications))).Methods(http.MethodGet)

	return router
}
