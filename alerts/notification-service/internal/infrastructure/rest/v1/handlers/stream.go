package handlers

import (
	"fmt"
	"net/http"
	"notification-serivce/internal/infrastructure/server"
	"time"
)

func (h *RequestHandler) StreamNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	val := r.Context().Value("id")
	clientID, ok := val.(string)
	if !ok || clientID == "" {
		errorMessage := "Missing client ID in context"
		server.NewResponse(w, &errorMessage, http.StatusBadRequest, nil)
		return
	}

	clientChan, err := h.manager.Subscribe(clientID)
	if err != nil {
		errorMessage := "Failed to subscribe stream"
		server.NewResponse(w, &errorMessage, http.StatusBadRequest, nil)
		return
	}

	defer h.manager.Unsubscribe(clientID)

	flusher, ok := w.(http.Flusher)
	if !ok {
		errorMessage := "Unsupported streaming"
		server.NewResponse(w, &errorMessage, http.StatusBadRequest, nil)
		return
	}

	//keeping alive
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case notification := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", notification)
			flusher.Flush()
		}
	}
}
