package handlers

import (
	"fmt"
	"log"
	"net/http"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/config"
	"notification-serivce/internal/infrastructure/server"
	"notification-serivce/internal/infrastructure/sse"
	"os"
)

var HEADERS = map[string]string{
	"Content-Type":                     "text/event-stream",
	"Cache-Control":                    "no-cache",
	"Connection":                       "keep-alive",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Headers":     "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Requested-With",
	"X-Accel-Buffering":                "no",
	"Transfer-Encoding":                "chunked",
}

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

// StreamNotifications godoc
// @Summary Stream user notifications
// @Description Server-Sent Events (SSE) stream for sending notifications in real-time.
// @Tags notifications
// @Produce text/event-stream
// @Param token query string true "JWT token for auth"
// @Router /api/v1/notification/stream [get]
func (h *SSEHandler) StreamNotifications(w http.ResponseWriter, r *http.Request) {
	h.configureSSEHeaders(w)

	conn, err := h.prepareSSEConnection(w, r)
	if err != nil {
		log.Println(err)
		h.handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		conn.Cleanup()
	}()

	if err := conn.SendWelcomeMessage(); err != nil {
		log.Printf("Failed to send welcome message to client %s: %v", conn.GetClientID(), err)
		return
	}

	h.manager.FlushBufferedNotification(conn.GetClientID(), conn.GetChannel())

	log.Printf("SSE connection established for client: %s", conn.GetClientID())

	conn.HandleEvents()
}

func (h *SSEHandler) prepareSSEConnection(w http.ResponseWriter, r *http.Request) (sse.NotificationStreamConnectionInterface, error) {
	ctx := r.Context()
	clientID, err := ExtractClientID(r)
	if err != nil {
		return nil, err
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming unsupported")
	}

	clientChan, err := h.manager.Subscribe(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe stream: %w", err)
	}

	conn := sse.NewNotificationStreamConnection(sse.NotificationStreamConnectionConfig{
		ClientID:          clientID,
		Writer:            w,
		Flusher:           flusher,
		Channel:           clientChan,
		Context:           ctx,
		HeartbeatInterval: h.config.HeartbeatInterval,
		Manager:           h.manager,
	})

	return conn, nil
}

func (h *SSEHandler) configureSSEHeaders(w http.ResponseWriter) {
	origin := os.Getenv(config.FRONTEND_URL)
	if origin == "" {
		origin = h.config.FrontendURL
	}

	headers := HEADERS
	headers["Access-Control-Allow-Origin"] = origin

	for key, value := range headers {
		w.Header().Set(key, value)
	}

	log.Printf("SSE headers configured with origin: %s", origin)
}

func (h *SSEHandler) handleError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("SSE Error: %s", message)
	server.NewResponse(w, &message, statusCode, nil)
}
