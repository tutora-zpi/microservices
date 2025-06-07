package rest

import (
	"meeting-scheduler-service/internal/infrastructure/handlers"
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(meetingHandler handlers.ManageMeetingHandler) *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1/meeting").Subrouter()

	api.Handle("/start", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.Handler))))
	api.Handle("/end", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.Handler))))

	return router
}
