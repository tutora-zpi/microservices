package notificationmanager

import (
	"log"
	"notification-serivce/internal/domain/buffer"
	"notification-serivce/internal/domain/dto"
	"sync"
	"time"
)

type NotificationBuffer struct {
	mutex        sync.RWMutex
	buffers      map[string][]*buffer.BufferedNotification
	maxSize      int
	ttl          time.Duration
	acknowledged map[string]map[int64]bool // clientID -> notificationID -> acknowledged
}

func NewNotificationBuffer(maxSize int, ttl time.Duration) *NotificationBuffer {
	nb := &NotificationBuffer{
		buffers:      make(map[string][]*buffer.BufferedNotification),
		acknowledged: make(map[string]map[int64]bool),
		maxSize:      maxSize,
		ttl:          ttl,
	}

	go nb.cleanupRoutine()
	return nb
}

func (nb *NotificationBuffer) AddNotification(dto dto.NotificationDTO) {
	nb.mutex.Lock()
	defer nb.mutex.Unlock()

	clientID := dto.Receiver.ID
	msg := buffer.NewBufferedNotification(dto)

	nb.buffers[clientID] = append(nb.buffers[clientID], msg)
	nb.ensureAcknowledgedMap(clientID)

	if len(nb.buffers[clientID]) > nb.maxSize {
		removed := nb.buffers[clientID][0]
		nb.buffers[clientID] = nb.buffers[clientID][1:]
		log.Printf("Removed oldest notification %d from client %s", removed.ID, clientID)
	}
}

func (nb *NotificationBuffer) GetBufferedNotifications(clientID string) []*buffer.BufferedNotification {
	nb.mutex.RLock()
	defer nb.mutex.RUnlock()

	now := time.Now()
	valid := nb.filterValidNotifications(clientID, now)

	log.Printf("Returning %d valid notifications for client %s", len(valid), clientID)
	return valid
}

func (nb *NotificationBuffer) AcknowledgeNotification(clientID string, notificationID int64) {
	nb.mutex.Lock()
	defer nb.mutex.Unlock()
	nb.ensureAcknowledgedMap(clientID)
	nb.acknowledged[clientID][notificationID] = true
}

func (nb *NotificationBuffer) MarkAllAsAcknowledged(clientID string) {
	nb.mutex.Lock()
	defer nb.mutex.Unlock()
	nb.ensureAcknowledgedMap(clientID)

	for _, msg := range nb.buffers[clientID] {
		nb.acknowledged[clientID][msg.ID] = true
	}
}

func (nb *NotificationBuffer) ClearExpiredNotifications(clientIDs ...string) {
	nb.mutex.Lock()
	defer nb.mutex.Unlock()

	now := time.Now()

	if len(clientIDs) > 0 {
		for _, id := range clientIDs {
			nb.clearForClient(id, now)
		}
	} else {
		for clientID := range nb.buffers {
			nb.clearForClient(clientID, now)
		}
	}
}

func (nb *NotificationBuffer) clearForClient(clientID string, now time.Time) {
	notifications := nb.buffers[clientID]
	if len(notifications) == 0 {
		return
	}

	var valid []*buffer.BufferedNotification
	for _, msg := range notifications {
		if msg.Age(&now) < nb.ttl {
			valid = append(valid, msg)
		} else {
			if nb.acknowledged[clientID] != nil {
				delete(nb.acknowledged[clientID], msg.ID)
			}
			log.Printf("Expired notification %d removed for client %s", msg.ID, clientID)
		}
	}

	if len(valid) == 0 {
		delete(nb.buffers, clientID)
		delete(nb.acknowledged, clientID)
	} else {
		nb.buffers[clientID] = valid
	}
}

func (nb *NotificationBuffer) ensureAcknowledgedMap(clientID string) {
	if nb.acknowledged[clientID] == nil {
		nb.acknowledged[clientID] = make(map[int64]bool)
	}
}

func (nb *NotificationBuffer) isAcknowledged(clientID string, notificationID int64) bool {
	return nb.acknowledged[clientID] != nil && nb.acknowledged[clientID][notificationID]
}

func (nb *NotificationBuffer) filterValidNotifications(clientID string, now time.Time) []*buffer.BufferedNotification {
	var valid []*buffer.BufferedNotification
	for _, msg := range nb.buffers[clientID] {
		if msg.Age(&now) < nb.ttl && !nb.isAcknowledged(clientID, msg.ID) {
			valid = append(valid, msg)
		}
	}
	return valid
}

func (nb *NotificationBuffer) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		nb.ClearExpiredNotifications()
	}
}
