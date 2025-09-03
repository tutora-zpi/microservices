package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/config"
	"notification-serivce/internal/infrastructure/server"
	"notification-serivce/internal/infrastructure/sse"
	"os"
)

type SSEHandler struct {
	manager interfaces.NotificationManager
	config  *sse.SSEConfig
}

func NewSSEHandler(manager interfaces.NotificationManager, cfg *sse.SSEConfig) *SSEHandler {
	if cfg == nil {
		cfg = &sse.SSEConfig{
			HeartbeatInterval: DefaultHeartbeatInterval,
			FrontendURL:       "*",
			MaxConnections:    1000,
		}
	}

	return &SSEHandler{
		manager: manager,
		config:  cfg,
	}
}

func (h *SSEHandler) prepareSSEConnection(w http.ResponseWriter, r *http.Request) (*sse.SSEConnection, context.CancelFunc, error) {
	clientID, err := ExtractClientID(r)
	if err != nil {
		return nil, nil, err
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, nil, fmt.Errorf("streaming unsupported")
	}

	clientChan, cancel, err := h.manager.Subscribe(clientID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe stream: %w", err)
	}

	conn := sse.NewSSEConnection(sse.ConnectionConfig{
		ClientID:          clientID,
		Writer:            w,
		Flusher:           flusher,
		Channel:           clientChan,
		Context:           r.Context(),
		HeartbeatInterval: h.config.HeartbeatInterval,
		Manager:           h.manager,
	})

	return conn, cancel, nil
}

func (h *SSEHandler) configureSSEHeaders(w http.ResponseWriter) {
	origin := os.Getenv(config.FRONTEND_URL)
	if origin == "" {
		origin = h.config.FrontendURL
	}

	headers := map[string]string{
		"Content-Type":                     "text/event-stream",
		"Cache-Control":                    "no-cache",
		"Connection":                       "keep-alive",
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Requested-With",
		"X-Accel-Buffering":                "no",
		"Transfer-Encoding":                "chunked",
	}

	for key, value := range headers {
		w.Header().Set(key, value)
	}

	log.Printf("SSE headers configured with origin: %s", origin)
}

func (h *SSEHandler) handleError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("SSE Error: %s", message)
	server.NewResponse(w, &message, statusCode, nil)
}
