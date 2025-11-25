package sse

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/pkg"
	"time"
)

type NotificationStreamConnectionInterface interface {
	HandleEvents()
	Cleanup()

	GetClientID() string
	GetChannel() chan []byte

	SendWelcomeMessage() error
	SendSSEComment(message string) error
	SendSSEEvent(eventType string, data []byte) error
	SendRetry(retryAfter int) error
}

type NotificationStreamConnection struct {
	ClientID        string
	Writer          http.ResponseWriter
	Flusher         http.Flusher
	Channel         chan []byte
	Context         context.Context
	HeartbeatTicker *time.Ticker
	Manager         interfaces.NotificationManager

	NotificationsSent int
	HeartbeatsSent    int
	ConnectionTime    time.Time
	lastActivity      time.Time
	isHealthy         bool
}

func (conn *NotificationStreamConnection) GetChannel() chan []byte {
	return conn.Channel
}

func (conn *NotificationStreamConnection) GetClientID() string {
	return conn.ClientID
}

func NewNotificationStreamConnection(config NotificationStreamConnectionConfig) NotificationStreamConnectionInterface {
	now := time.Now()
	return &NotificationStreamConnection{
		ClientID:        config.ClientID,
		Writer:          config.Writer,
		Flusher:         config.Flusher,
		Channel:         config.Channel,
		Context:         config.Context,
		HeartbeatTicker: time.NewTicker(config.HeartbeatInterval),
		Manager:         config.Manager,
		ConnectionTime:  now,
		lastActivity:    now,
		isHealthy:       true,
	}
}

func (conn *NotificationStreamConnection) HandleEvents() {
	_ = conn.SendRetry(5)

	_ = conn.SendSSEComment("Connection ready, checking for buffered notifications...")

	conn.Flusher.Flush()

	connectionTimeout := time.NewTimer(5 * time.Minute)
	defer connectionTimeout.Stop()

	for {
		select {
		case <-conn.Context.Done():
			log.Printf("SSE connection closed for client: %s (context done) - Stats: %d notifications, %d heartbeats, duration: %v",
				conn.ClientID, conn.NotificationsSent, conn.HeartbeatsSent, time.Since(conn.ConnectionTime))
			return

		case notification, ok := <-conn.Channel:
			if !ok {
				log.Printf("Notification channel closed for client: %s", conn.ClientID)
				return
			}

			if err := conn.SendSSEEvent("notification", notification); err != nil {
				log.Printf("Failed to send notification to client %s: %v", conn.ClientID, err)
				return
			}

			log.Printf("Sent notification to client %s (total: %d)",
				conn.ClientID, conn.NotificationsSent)

			connectionTimeout.Reset(5 * time.Minute)

		case <-conn.HeartbeatTicker.C:
			if !conn.checkConnectionHealth() {
				log.Printf("Connection health check failed for client %s", conn.ClientID)
				return
			}

			if err := conn.sendHeartbeat(); err != nil {
				log.Printf("Failed to send heartbeat to client %s: %v", conn.ClientID, err)
				return
			}

			connectionTimeout.Reset(5 * time.Minute)

		case <-connectionTimeout.C:
			log.Printf("Connection timeout for client %s (no activity for 5 minutes)", conn.ClientID)
			return
		}
	}
}

func (conn *NotificationStreamConnection) SendWelcomeMessage() error {
	return conn.SendSSEComment(fmt.Sprintf("SSE connection established for client %s at %s",
		conn.ClientID, conn.ConnectionTime.Format(time.RFC3339)))
}

func (conn *NotificationStreamConnection) SendSSEComment(message string) error {
	if !conn.isHealthy {
		return fmt.Errorf("connection is not healthy")
	}

	_, err := fmt.Fprintf(conn.Writer, "data: %s\n\n", message)
	if err != nil {
		conn.isHealthy = false
		log.Println("An error occured during sending ping")
		return fmt.Errorf("failed to write SSE comment: %w", err)
	}
	conn.Flusher.Flush()
	conn.lastActivity = time.Now()
	return nil
}

func (conn *NotificationStreamConnection) SendSSEEvent(eventType string, data []byte) error {
	if !conn.isHealthy {
		return fmt.Errorf("connection is not healthy")
	}

	if len(data) == 0 {
		log.Printf("Warning: Sending empty data to client %s", conn.ClientID)
	}

	timestamp := pkg.GenerateTimestamp()

	message := fmt.Sprintf("id: %d\nevent: %s\ndata: %s\n\n", timestamp, eventType, data)

	log.Printf("Sending SSE event to client %s: type=%s, size=%d bytes",
		conn.ClientID, eventType, len(message))

	_, err := conn.Writer.Write([]byte(message))
	if err != nil {
		conn.isHealthy = false
		log.Printf("Write error for client %s: %v", conn.ClientID, err)
		return fmt.Errorf("failed to write SSE event: %w", err)
	}

	conn.Flusher.Flush()
	conn.NotificationsSent++
	conn.lastActivity = time.Now()
	return nil
}

func (conn *NotificationStreamConnection) SendRetry(retryAfter int) error {
	_, err := fmt.Fprintf(conn.Writer, "retry: %d\n\n", retryAfter*1000) // milliseconds
	if err != nil {
		return fmt.Errorf("failed to write SSE retry: %w", err)
	}
	conn.Flusher.Flush()
	return nil
}

func (conn *NotificationStreamConnection) Cleanup() {
	if conn.HeartbeatTicker != nil {
		conn.HeartbeatTicker.Stop()
	}

	conn.Manager.Unsubscribe(conn.ClientID)
	log.Printf("Immediate unsubscribe for client: %s", conn.ClientID)

	stats := conn.GetStats()
	log.Printf("SSE connection cleaned up for client: %s - Final stats: %+v", conn.ClientID, stats)
}

func (conn *NotificationStreamConnection) sendHeartbeat() error {
	log.Println("Sending heartbeat")
	err := conn.SendSSEComment(fmt.Sprintf("heartbeat-%d", time.Now().Unix()))
	if err == nil {
		conn.HeartbeatsSent++
	}
	return err
}

func (conn *NotificationStreamConnection) checkConnectionHealth() bool {
	if time.Since(conn.lastActivity) > 5*time.Minute {
		log.Printf("Connection for client %s appears stale (last activity: %v ago)",
			conn.ClientID, time.Since(conn.lastActivity))
		return false
	}
	return conn.isHealthy
}
