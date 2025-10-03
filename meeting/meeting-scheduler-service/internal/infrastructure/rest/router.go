package rest

import (
	"meeting-scheduler-service/internal/infrastructure/middleware"
	"meeting-scheduler-service/internal/infrastructure/rest/v1/handlers"
	"net/http"

	_ "meeting-scheduler-service/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(meetingHandler handlers.ManageMeetingHandler) *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	api.Handle("/docs", http.RedirectHandler("/api/v1/docs/", http.StatusSeeOther))

	api.Handle("/meeting/start", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.StartMeeting)))).Methods(http.MethodPost, http.MethodPut)
	api.Handle("/meeting/end", middleware.IsAuth(middleware.Validate(http.HandlerFunc(meetingHandler.EndMeeting)))).Methods(http.MethodPost, http.MethodDelete)
	api.Handle("/meeting/{class_id}", middleware.IsAuth(http.HandlerFunc(meetingHandler.GetActiveMeeting))).Methods(http.MethodGet)

	return router
}
