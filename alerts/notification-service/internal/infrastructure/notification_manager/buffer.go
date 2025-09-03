package notificationmanager

import (
	"log"
	"notification-serivce/internal/domain/buffer"
	"notification-serivce/internal/domain/dto"
	"sync"
	"time"
)

type NotificationBuffer struct {
	mu      sync.RWMutex
	buffers map[string][]*buffer.BufferedNotification
	maxSize int
	ttl     time.Duration

	// clientID -> notificationID -> acknowledged
	acknowledged map[string]map[string]bool
}

func NewNotificationBuffer(maxSize int, ttl time.Duration) *NotificationBuffer {
	buffer := &NotificationBuffer{
		buffers:      make(map[string][]*buffer.BufferedNotification),
		acknowledged: make(map[string]map[string]bool),
		maxSize:      maxSize,
		ttl:          ttl,
	}

	go buffer.cleanupRoutine()

	return buffer
}

func (mb *NotificationBuffer) AddNotification(dto dto.NotificationDTO) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	log.Println("Appending notification to the buffer")

	clientID := dto.ReceiverID

	msg := buffer.NewBufferedNotification(dto)

	mb.buffers[clientID] = append(mb.buffers[clientID], msg)

	currentSize := len(mb.buffers[clientID])

	if currentSize > mb.maxSize {
		removedNotification := mb.buffers[clientID][0]

		mb.buffers[clientID] = mb.buffers[clientID][1:]

		log.Printf("Removed oldest notification %s from client %s", removedNotification.ID, clientID)
	}
}

func (mb *NotificationBuffer) GetBufferedNotifications(clientID string) []*buffer.BufferedNotification {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	notifications := mb.buffers[clientID]
	now := time.Now()

	log.Printf("Get notifications for client %s - total: %d", clientID, len(notifications))

	var validNotifications []*buffer.BufferedNotification
	for _, msg := range notifications {

		age := now.Sub(msg.Timestamp)

		if age < mb.ttl {

			if mb.acknowledged[clientID] == nil || !mb.acknowledged[clientID][msg.ID] {
				validNotifications = append(validNotifications, msg)
				log.Printf("Valid notification %s (age: %v)", msg.ID, age)

			} else {
				log.Printf("Notification %s already acknowledged", msg.ID)
			}

		} else {
			log.Printf("Notification %s expired (age: %v, TTL: %v)", msg.ID, age, mb.ttl)
		}
	}

	log.Printf("Returning %d valid notifications for client %s", len(validNotifications), clientID)
	return validNotifications
}

func (mb *NotificationBuffer) AcknowledgeNotification(clientID string, notificationID string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if mb.acknowledged[clientID] == nil {
		mb.acknowledged[clientID] = make(map[string]bool)
	}
	mb.acknowledged[clientID][notificationID] = true
}

func (mb *NotificationBuffer) MarkAllAsAcknowledged(clientID string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if mb.acknowledged[clientID] == nil {
		mb.acknowledged[clientID] = make(map[string]bool)
	}

	for _, msg := range mb.buffers[clientID] {
		mb.acknowledged[clientID][msg.ID] = true
	}
}

func (mb *NotificationBuffer) ClearExpiredNotifications(clientID ...string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	now := time.Now()

	if len(clientID) > 0 {
		for _, id := range clientID {
			mb.clearForClient(id, now)
		}
	} else {
		for clientID := range mb.buffers {
			mb.clearForClient(clientID, now)
		}
	}
}

func (mb *NotificationBuffer) clearForClient(clientID string, now time.Time) {
	notifications := mb.buffers[clientID]
	if len(notifications) == 0 {
		return
	}

	var validNotifications []*buffer.BufferedNotification

	for _, msg := range notifications {
		if now.Sub(msg.Timestamp) < mb.ttl {
			validNotifications = append(validNotifications, msg)
		} else {
			if mb.acknowledged[clientID] != nil {
				delete(mb.acknowledged[clientID], msg.ID)
			}
		}
	}

	if len(validNotifications) == 0 {
		delete(mb.buffers, clientID)
		delete(mb.acknowledged, clientID)
	} else {
		mb.buffers[clientID] = validNotifications
	}
}

func (mb *NotificationBuffer) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mb.ClearExpiredNotifications()
	}
}
