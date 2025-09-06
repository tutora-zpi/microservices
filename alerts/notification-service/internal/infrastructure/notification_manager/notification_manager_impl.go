package notificationmanager

import (
	"context"
	"fmt"
	"log"
	"notification-serivce/internal/app/interfaces"
	"notification-serivce/internal/domain/buffer"
	"notification-serivce/internal/domain/dto"
	"sync"
	"time"
)

type notificationManagerImpl struct {
	mu                 sync.RWMutex
	clients            map[string]chan []byte
	notificationBuffer *NotificationBuffer
	bufferingEnabled   bool
	connectionTracker  map[string]time.Time
}

func NewManager() interfaces.NotificationManager {
	manager := &notificationManagerImpl{
		clients:            make(map[string]chan []byte),
		connectionTracker:  make(map[string]time.Time),
		bufferingEnabled:   false,
		notificationBuffer: nil,
	}
	log.Printf("NotificationManager created - buffering disabled by default")
	return manager
}

func (m *notificationManagerImpl) EnableBuffering(maxSize int, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.notificationBuffer = NewNotificationBuffer(maxSize, ttl)
	m.bufferingEnabled = true
	log.Printf("Buffering ENABLED: maxSize=%d, ttl=%v", maxSize, ttl)
}

func (m *notificationManagerImpl) Subscribe(clientID string) (chan []byte, context.CancelFunc, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	lastConnection := m.connectionTracker[clientID]
	m.connectionTracker[clientID] = time.Now()

	oldChan, exists := m.clients[clientID]

	clientChan := make(chan []byte, 200)
	_, cancel := context.WithCancel(context.Background())

	m.clients[clientID] = clientChan

	log.Printf("Client %s subscribed (total clients: %d)", clientID, len(m.clients))

	if exists && oldChan != nil {
		close(oldChan)
		log.Printf("Closed old channel for reconnecting client %s", clientID)
	}

	if !lastConnection.IsZero() && time.Since(lastConnection) < 2*time.Minute {
		log.Printf("Quick reconnect detected for client %s (last seen: %v ago)",
			clientID, time.Since(lastConnection))
	}

	return clientChan, cancel, nil
}

func (m *notificationManagerImpl) GetBufferedNotifications(clientID string) []*buffer.BufferedNotification {
	if m.bufferingEnabled && m.notificationBuffer != nil {
		return m.notificationBuffer.GetBufferedNotifications(clientID)
	}
	return nil
}

func (m *notificationManagerImpl) Push(dto dto.NotificationDTO) error {
	m.mu.RLock()
	clientChan, exists := m.clients[dto.Receiver.ID]
	clientCount := len(m.clients)
	m.mu.RUnlock()

	data := dto.JSON()
	log.Printf("Client %s, exists: %t, total clients: %d, buffering: %t",
		dto.Receiver.ID, exists, clientCount, m.bufferingEnabled)

	isActuallyConnected := exists && clientChan != nil

	if !isActuallyConnected {
		if m.bufferingEnabled && m.notificationBuffer != nil {
			m.notificationBuffer.AddNotification(dto)
			log.Println("Notification buffered")
		} else {
			log.Println("Buffering disabled")
		}
		return nil
	}

	select {
	case clientChan <- data:
		log.Printf("Notification sent to client %s", dto.Receiver.ID)
		return nil

	default:
		if m.bufferingEnabled && m.notificationBuffer != nil {
			m.notificationBuffer.AddNotification(dto)
			log.Printf("Channel full for client %s - notification buffered\n", dto.Receiver.ID)
		} else {
			log.Printf("Channel full for client %s - buffering disabled\n", dto.Receiver.ID)
		}
		return fmt.Errorf("client %s channel full", dto.Receiver.ID)
	}
}

func (m *notificationManagerImpl) IsClientConnected(clientID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clientChan, exists := m.clients[clientID]
	return exists && clientChan != nil
}

func (m *notificationManagerImpl) Unsubscribe(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if clientChan, exists := m.clients[clientID]; exists {
		close(clientChan)
		delete(m.clients, clientID)
		log.Printf("Client %s unsubscribed (total clients: %d)", clientID, len(m.clients))
	}
}

func (m *notificationManagerImpl) FlushBufferedNotification(clientID string, clientChan chan []byte) int {
	log.Println("Flushing buffered notifications")

	if !m.bufferingEnabled || m.notificationBuffer == nil {
		return 0
	}

	buffered := m.notificationBuffer.GetBufferedNotifications(clientID)
	sentCount := 0

	for _, msg := range buffered {
		select {
		case clientChan <- msg.Data:
			sentCount++
			m.notificationBuffer.AcknowledgeNotification(clientID, msg.ID)
		default:
			log.Printf("Channel full, stopping buffered notifications flush for client %s", clientID)
			return sentCount
		}
	}

	if sentCount > 0 {
		log.Printf("Flushed %d buffered notificationss to client %s", sentCount, clientID)
	}
	return sentCount
}

func (m *notificationManagerImpl) AcknowledgeNotification(clientID string, notificationID string) {
	if m.bufferingEnabled && m.notificationBuffer != nil {
		m.notificationBuffer.AcknowledgeNotification(clientID, notificationID)
	}
}
