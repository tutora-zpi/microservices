package sse

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"sync"
)

type SSEManager struct {
	subscribers map[string]chan []byte
	lock        sync.Mutex
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		subscribers: make(map[string]chan []byte),
	}
}

func (m *SSEManager) Subscribe(clientID string) (chan []byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	ch := make(chan []byte, 32)

	// only unique clientIDs
	if _, ok := m.subscribers[clientID]; !ok {
		m.subscribers[clientID] = ch
	}

	return ch, nil
}

func (m *SSEManager) Unsubscribe(clientID string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	ch, ok := m.subscribers[clientID]
	if !ok {
		return fmt.Errorf("client not found")
	}
	close(ch)
	delete(m.subscribers, clientID)
	return nil
}

func (m *SSEManager) Push(dto dto.NotificationDTO) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	data := dto.JSON()

	for _, ch := range m.subscribers {
		select {
		case ch <- data:
		default:
		}
	}
	return nil
}
