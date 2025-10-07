package v1

import (
	"net/http"
	_ "notification-serivce/docs"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/infrastructure/middleware"
	"notification-serivce/internal/infrastructure/rest/v1/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(manager interfaces.NotificationManager, service interfaces.NotificationSerivce) *mux.Router {
	router := mux.NewRouter()
	sseHandler := handlers.NewSSEHandler(manager, nil)
	httpHandler := handlers.NewHTTPHandler(service)

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	api.Handle("/docs", http.RedirectHandler("/api/v1/docs/", http.StatusSeeOther))

	api.Handle("/stream", middleware.IsAuth(http.HandlerFunc(sseHandler.StreamNotifications))).Methods(http.MethodGet)

	notifcation := api.PathPrefix("/notification").Subrouter()
	// notifcation.Handle("", middleware.IsAuth(http.HandlerFunc(httpHandler.FetchNotifications))).Methods(http.MethodGet)
	// notifcation.Handle("", middleware.IsAuth(http.HandlerFunc(httpHandler.DeleteNotifications))).Methods(http.MethodDelete)
	notifcation.Handle("", http.HandlerFunc(httpHandler.FetchNotifications)).Methods(http.MethodGet)
	notifcation.Handle("", http.HandlerFunc(httpHandler.DeleteNotifications)).Methods(http.MethodDelete)

	return router
}
