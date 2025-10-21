package ws

import (
	"log"
	"sync"
	"ws-gateway/internal/app/interfaces"

	"github.com/gorilla/websocket"
)

type Room struct {
	mu      sync.RWMutex
	members map[string]interfaces.Client
}

type hub struct {
	globalMembers sync.Map // map[userid]client -> set
	rooms         sync.Map // map[roomID]map[userid]client -> room id -> set
	close         sync.Once
}

// Close implements interfaces.HubManager.
func (h *hub) Close() {
	h.close.Do(func() {
		h.globalMembers.Range(func(userID, client any) bool {
			if c, ok := client.(interfaces.Client); ok {
				c.GetConnection().Close()
			}
			return true
		})

		h.rooms.Range(func(roomID, members any) bool {
			if r, ok := members.(*Room); ok {
				r.mu.RLock()
				for _, client := range r.members {
					client.GetConnection().Close()
				}
				r.mu.RUnlock()
			}

			return true
		})

		h.globalMembers.Clear()
		h.rooms.Clear()

		log.Println("WSocket closed successfully.")
	})
}

// EmitToClient implements interfaces.HubManager.
func (h *hub) EmitToClient(clientID string, payloads [][]byte) {
	clientInterface, ok := h.globalMembers.Load(clientID)
	if !ok {
		log.Printf("cliendID %s not found", clientID)
		return
	}

	client := clientInterface.(interfaces.Client)

	for _, payload := range payloads {
		err := client.GetConnection().WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Printf("An error occurred during writing an message: %v", err)
		}
	}
}

// EmitToClientInRoom implements interfaces.HubManager.
func (h *hub) EmitToClientInRoom(roomID string, clientID string, payloads [][]byte) {
	roomInterface, ok := h.rooms.Load(roomID)
	if !ok {
		log.Printf("roomID %s not found", roomID)
		return
	}

	room := roomInterface.(*Room)

	room.mu.RLock()
	defer room.mu.RUnlock()

	client := room.members[clientID]
	for _, payload := range payloads {
		err := client.GetConnection().WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Printf("An error occurred during writing an message: %v", err)
		}
	}
}

// RemoveRoomMember implements interfaces.HubManager.
func (h *hub) RemoveRoomMember(roomID string, client interfaces.Client) (roomUsers []string) {
	roomInterface, ok := h.rooms.Load(roomID)
	if !ok {
		log.Printf("roomID %s not found", roomID)
		return
	}

	room := roomInterface.(*Room)

	room.mu.Lock()
	defer room.mu.Unlock()

	delete(room.members, client.ID())

	if len(room.members) == 0 {
		log.Printf("Removing room %s - no members", roomID)
		h.rooms.Delete(roomID)
		return []string{}
	}

	keys := make([]string, 0, len(room.members))
	for k := range room.members {
		keys = append(keys, k)
	}

	return keys
}

// RemoveGlobalMember implements interfaces.HubManager.
func (h *hub) RemoveGlobalMember(client interfaces.Client) {
	log.Printf("Removing new user: %s", client.ID())
	h.globalMembers.Delete(client.ID())
}

// AddGlobalMember implements interfaces.HubManager.
func (h *hub) AddGlobalMember(client interfaces.Client) {
	log.Printf("Adding new user: %s", client.ID())
	h.globalMembers.Store(client.ID(), client)
}

// AddRoomMember implements interfaces.HubManager.
func (h *hub) AddRoomMember(roomID string, c interfaces.Client) (roomUsers []string) {
	roomInterface, _ := h.rooms.LoadOrStore(roomID, &Room{members: make(map[string]interfaces.Client)})
	room := roomInterface.(*Room)

	room.mu.Lock()
	defer room.mu.Unlock()
	room.members[c.ID()] = c

	keys := make([]string, 0, len(room.members))
	for k := range room.members {
		keys = append(keys, k)
	}
	return keys
}

// EmitGlobal implements interfaces.HubManager.
func (h *hub) EmitGlobal(payload []byte) {
	h.globalMembers.Range(func(_, value any) bool {
		client := value.(interfaces.Client)
		if err := client.GetConnection().WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("Error sending to global %s: %v", client.ID(), err)
		}
		return true
	})
}

// Emit implements interfaces.HubManager.
func (h *hub) Emit(roomID string, payload []byte, pred func(id string) bool) {
	log.Println("Emiiting")
	value, ok := h.rooms.Load(roomID)
	if !ok {
		log.Printf("No room with id: %s", roomID)
		return
	}

	room := value.(*Room)

	room.mu.RLock()
	defer room.mu.RUnlock()

	for _, client := range room.members {
		if pred(client.ID()) {
			if err := client.GetConnection().WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Printf("Error sending to %s: %v", client.ID(), err)
			} else {
				log.Printf("Successfully emitted: %s", string(payload))
			}
		}
	}
}

func NewHub() interfaces.HubManager {
	return &hub{}
}
