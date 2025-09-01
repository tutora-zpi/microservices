package handlers

import "notification-serivce/internal/infrastructure/sse"

type RequestHandler struct {
	manager *sse.SSEManager
}

func NewRequestHandler(manager *sse.SSEManager) *RequestHandler {
	return &RequestHandler{
		manager: manager,
	}
}
