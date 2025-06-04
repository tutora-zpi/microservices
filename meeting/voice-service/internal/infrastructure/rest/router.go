package rest

import (
	"voice-service/internal/infrastructure/config"

	"github.com/gorilla/mux"
)

func NewRouter(i *config.Incjectable) *mux.Router {
	router := &mux.Router{}

	router.HandleFunc("/ws", i.Gateway.Handle)

	// router.HandleFunc("/health")

	return router
}
