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

type connectionInfo struct {
	channel      chan []byte
	createdAt    time.Time
	ctx          context.Context
	connectionID string
}

type notificationManagerImpl struct {
	mutex              sync.RWMutex
	clients            map[string][]*connectionInfo
	notificationBuffer *NotificationBuffer
	bufferingEnabled   bool
	connectionTracker  map[string]time.Time
}

func NewManager() interfaces.NotificationManager {
	manager := &notificationManagerImpl{
		clients:           make(map[string][]*connectionInfo),
		connectionTracker: make(map[string]time.Time),
		bufferingEnabled:  false,
	}
	return manager
}

func (m *notificationManagerImpl) EnableBuffering(maxSize int, ttl time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.notificationBuffer = NewNotificationBuffer(maxSize, ttl)
	m.bufferingEnabled = true
	log.Printf("Buffering ENABLED: maxSize=%d, ttl=%v", maxSize, ttl)
}

func (m *notificationManagerImpl) Subscribe(ctx context.Context, clientID string) (chan []byte, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	connectionID := fmt.Sprintf("%s-%d", clientID, time.Now().UnixNano())

	clientChan := make(chan []byte, 200)

	newConn := &connectionInfo{
		channel:      clientChan,
		createdAt:    time.Now(),
		ctx:          ctx,
		connectionID: connectionID,
	}

	m.clients[clientID] = append(m.clients[clientID], newConn)

	go m.monitorConnection(ctx, clientID, connectionID, clientChan)

	log.Printf("Client %s subscribed with connection %s (total connections for user: %d, all clients: %d)",
		clientID, connectionID, len(m.clients[clientID]), m.getTotalConnections())

	return clientChan, nil
}

func (m *notificationManagerImpl) getTotalConnections() int {
	total := 0
	for _, conns := range m.clients {
		total += len(conns)
	}
	return total
}

func (m *notificationManagerImpl) monitorConnection(ctx context.Context, clientID string, connectionID string, clientChan chan []byte) {
	<-ctx.Done()

	m.mutex.Lock()
	connections := m.clients[clientID]
	newConnections := make([]*connectionInfo, 0, len(connections))

	for _, conn := range connections {
		if conn.connectionID != connectionID {
			newConnections = append(newConnections, conn)
		}
	}

	if len(newConnections) == 0 {
		delete(m.clients, clientID)
		m.connectionTracker[clientID] = time.Now()
	} else {
		m.clients[clientID] = newConnections
	}

	remainingConns := len(newConnections)
	m.mutex.Unlock()

	log.Printf("Connection %s for client %s ended (remaining connections: %d)",
		connectionID, clientID, remainingConns)
}

func (m *notificationManagerImpl) Unsubscribe(clientID string) {
	m.mutex.Lock()
	connections, existed := m.clients[clientID]
	if existed {
		delete(m.clients, clientID)
		for _, conn := range connections {
			if conn.channel != nil {
				close(conn.channel)
			}
		}
	}
	m.connectionTracker[clientID] = time.Now()
	m.mutex.Unlock()

	if existed {
		log.Printf("Client %s unsubscribed (%d connections closed)", clientID, len(connections))
	}
}

func (m *notificationManagerImpl) Push(dto dto.NotificationDTO) error {
	m.mutex.RLock()
	connections := m.clients[dto.Receiver.ID]
	clientCount := len(connections)
	m.mutex.RUnlock()

	data := dto.JSON()

	if clientCount == 0 {
		log.Printf("Client %s has no active connections, buffering notification", dto.Receiver.ID)
		m.bufferIfEnabled(dto)
		return nil
	}

	log.Printf("Pushing notification to client %s (%d active connections)",
		dto.Receiver.ID, clientCount)

	successCount := 0
	for _, conn := range connections {
		select {
		case conn.channel <- data:
			successCount++
		default:
			log.Printf("Channel full for connection %s", conn.connectionID)
		}
	}

	if successCount > 0 {
		log.Printf("Notification sent to %d/%d connections for client %s",
			successCount, clientCount, dto.Receiver.ID)
		return nil
	}

	m.bufferIfEnabled(dto)
	return fmt.Errorf("all channels full for client %s", dto.Receiver.ID)
}

func (m *notificationManagerImpl) bufferIfEnabled(dto dto.NotificationDTO) {
	if m.bufferingEnabled && m.notificationBuffer != nil {
		m.notificationBuffer.AddNotification(dto)
		log.Println("Notification buffered")
	} else {
		log.Println("Buffering disabled")
	}
}

func (m *notificationManagerImpl) IsClientConnected(clientID string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	connections := m.clients[clientID]
	return len(connections) > 0
}

func (m *notificationManagerImpl) GetBufferedNotifications(clientID string) []*buffer.BufferedNotification {
	if m.bufferingEnabled && m.notificationBuffer != nil {
		return m.notificationBuffer.GetBufferedNotifications(clientID)
	}
	return nil
}

func (m *notificationManagerImpl) FlushBufferedNotification(clientID string, clientChan chan []byte) int {
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
		log.Printf("Flushed %d buffered notifications to client %s", sentCount, clientID)
	}
	return sentCount
}

func (m *notificationManagerImpl) AcknowledgeNotification(clientID string, notificationID int64) {
	if m.bufferingEnabled && m.notificationBuffer != nil {
		m.notificationBuffer.AcknowledgeNotification(clientID, notificationID)
	}
}
