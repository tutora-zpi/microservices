package ws

import (
	"log"
	"signaling-service/internal/app/interfaces"
	"sync"
)

type Meeting struct {
	mu      sync.RWMutex
	members map[string]interfaces.Client
}

type hub struct {
	globalMembers sync.Map // map[userid]client -> set
	meetings      sync.Map // map[meetingid]map][userid]client -> room id -> set
}

// RemoveMeetingMemeber implements interfaces.HubManager.
func (h *hub) RemoveMeetingMemeber(meetingID string, client interfaces.Client) {
	meetingInterface, _ := h.meetings.Load(meetingID)
	meeting := meetingInterface.(*Meeting)

	meeting.mu.Lock()
	defer meeting.mu.Unlock()
	delete(meeting.members, client.ID())
}

// RemoveGlobalMember implements interfaces.HubManager.
func (h *hub) RemoveGlobalMember(client interfaces.Client) {
	log.Printf("Removing new user: %s", client.ID())
	h.globalMembers.Delete(client.ID())
}

func (h *hub) AddGlobalMember(client interfaces.Client) {
	log.Printf("Adding new user: %s", client.ID())
	h.globalMembers.Store(client.ID(), client)
}

func (h *hub) AddMeetingMember(meetingID string, c interfaces.Client) {
	meetingInterface, _ := h.meetings.LoadOrStore(meetingID, &Meeting{members: make(map[string]interfaces.Client)})
	meeting := meetingInterface.(*Meeting)

	meeting.mu.Lock()
	defer meeting.mu.Unlock()
	meeting.members[c.ID()] = c
}

func (h *hub) EmitGlobal(messageType int, payload []byte) {
	h.globalMembers.Range(func(_, value any) bool {
		client := value.(interfaces.Client)
		if err := client.GetConnection().WriteMessage(messageType, payload); err != nil {
			log.Printf("Error sending to global %s: %v", client.ID(), err)
		}
		return true
	})
}

func (h *hub) Emit(meetingID string, messageType int, payload []byte, pred func(id string) bool) {
	value, ok := h.meetings.Load(meetingID)
	if !ok {
		return
	}

	meeting := value.(*Meeting)

	meeting.mu.RLock()
	defer meeting.mu.RUnlock()

	for _, client := range meeting.members {
		if pred(client.ID()) {
			if err := client.GetConnection().WriteMessage(messageType, payload); err != nil {
				log.Printf("Error sending to %s: %v", client.ID(), err)
			}
		}
	}
}

func NewHub() interfaces.HubManager {
	return &hub{}
}
