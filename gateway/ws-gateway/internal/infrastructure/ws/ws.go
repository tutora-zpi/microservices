package ws

import (
	"log"
	"sync"
	"time"

	"ws-gateway/internal/app/interfaces"
)

type Room struct {
	mu      sync.RWMutex
	members map[string]interfaces.Client
}

type hub struct {
	globalMembers sync.Map
	rooms         sync.Map
	register      chan interfaces.Client
	unregister    chan interfaces.Client
	closing       chan struct{}
	once          sync.Once
}

// RemoveBotFromRoom implements interfaces.HubManager.
func (h *hub) RemoveRoomMemberByID(roomID string, botID string) (roomUsers []string) {
	v, ok := h.rooms.Load(roomID)
	if !ok {
		return h.roomIDs(roomID)
	}

	room := v.(*Room)
	room.mu.Lock()
	delete(room.members, botID)
	empty := len(room.members) == 0
	room.mu.Unlock()

	if empty {
		h.rooms.Delete(roomID)
	}
	return h.roomIDs(roomID)

}

func NewHub() interfaces.HubManager {
	h := &hub{
		register:   make(chan interfaces.Client, 32),
		unregister: make(chan interfaces.Client, 32),
		closing:    make(chan struct{}),
	}
	go h.run()
	return h
}

func (h *hub) GetUsersFromRoomID(roomID string) []string {
	v, ok := h.rooms.Load(roomID)
	if !ok {
		return nil
	}

	room := v.(*Room)
	room.mu.RLock()
	defer room.mu.RUnlock()
	ids := make([]string, 0, len(room.members))
	for k := range room.members {
		ids = append(ids, k)
	}

	return ids
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.globalMembers.Store(c.ID(), c)
		case c := <-h.unregister:
			h.globalMembers.Delete(c.ID())

			h.rooms.Range(func(_ any, v any) bool {
				r := v.(*Room)
				r.mu.Lock()
				delete(r.members, c.ID())
				r.mu.Unlock()
				return true
			})

		case <-h.closing:
			return
		}
	}
}

func (h *hub) AddGlobalMember(c interfaces.Client) {
	h.register <- c
}

func (h *hub) RemoveGlobalMember(c interfaces.Client) {
	h.unregister <- c
}

func (h *hub) AddRoomMember(roomID string, c interfaces.Client) []string {
	r, _ := h.rooms.LoadOrStore(roomID, &Room{members: make(map[string]interfaces.Client)})
	room := r.(*Room)

	room.mu.Lock()
	room.members[c.ID()] = c
	room.mu.Unlock()

	return h.roomIDs(roomID)
}

func (h *hub) RemoveRoomMember(roomID string, c interfaces.Client) []string {
	return h.RemoveRoomMemberByID(roomID, c.ID())
}

func (h *hub) roomIDs(roomID string) []string {
	v, ok := h.rooms.Load(roomID)
	if !ok {
		return nil
	}

	room := v.(*Room)
	room.mu.RLock()
	defer room.mu.RUnlock()

	ids := make([]string, 0, len(room.members))
	for id := range room.members {
		ids = append(ids, id)
	}
	return ids
}

func (h *hub) Emit(roomID string, payload []byte, pred func(id string) bool) {
	v, ok := h.rooms.Load(roomID)
	if !ok {
		return
	}

	room := v.(*Room)
	room.mu.RLock()
	defer room.mu.RUnlock()

	for _, c := range room.members {
		if pred == nil || pred(c.ID()) {
			h.sendSafe(c, payload)
		}
	}
}

func (h *hub) EmitGlobal(payload []byte) {
	h.globalMembers.Range(func(_ any, v any) bool {
		h.sendSafe(v.(interfaces.Client), payload)
		return true
	})
}

func (h *hub) EmitToClient(clientID string, payloads [][]byte) {
	v, _ := h.globalMembers.Load(clientID)
	if v == nil {
		return
	}
	c := v.(interfaces.Client)

	for _, p := range payloads {
		h.sendSafe(c, p)
	}
}

func (h *hub) EmitToClientInRoom(roomID, clientID string, payloads [][]byte) {
	v, ok := h.rooms.Load(roomID)
	if !ok {
		return
	}

	room := v.(*Room)

	room.mu.RLock()
	c, ok := room.members[clientID]
	room.mu.RUnlock()

	if !ok {
		return
	}

	for _, p := range payloads {
		h.sendSafe(c, p)
	}
}

func (h *hub) sendSafe(c interfaces.Client, payload []byte) {
	cli, ok := c.(*clientImpl)
	if !ok {
		return
	}

	select {
	case <-cli.done:
		return
	default:
	}

	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case cli.send <- payload:
	case <-timer.C:
		log.Printf("[TIMEOUT] → %s", c.ID())
	case <-cli.done:
		return
	}
}

func (h *hub) Close() {
	h.once.Do(func() {
		close(h.closing)
		h.globalMembers.Range(func(_ any, v any) bool {
			client, ok := v.(interfaces.Client)
			if ok {
				client.Close()
			}
			return ok
		})
	})
}
