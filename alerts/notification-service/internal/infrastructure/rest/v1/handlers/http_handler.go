package handlers

import "notification-serivce/internal/app/interfaces"

type HTTPHandler struct {
	bus interfaces.QueryBus
}

func NewHTTPHandler(bus interfaces.QueryBus) *HTTPHandler {
	return &HTTPHandler{bus: bus}
}
