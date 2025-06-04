package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"voice-service/internal/infrastructure/security"

	"github.com/gorilla/websocket"
)

// Query string ws://host:port/ws?token=...&room_id=...
type InitMessage struct {
	Token  string `json:"token"`
	RoomID string `json:"room_id"`
}

type gatewayImpl struct {
	upgrader websocket.Upgrader

	// roomID -> userID -> Client
	rooms map[string]map[string]Client
	mutex sync.Mutex

	jwtSecretKey []byte
}

func NewGateway(secret string) Gateway {
	return &gatewayImpl{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		rooms:        make(map[string]map[string]Client),
		jwtSecretKey: []byte(secret),
	}
}

func (g *gatewayImpl) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	_, initMsg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Init message read error:", err)
		conn.Close()
		return
	}

	var init InitMessage
	if err := json.Unmarshal(initMsg, &init); err != nil {
		log.Println("Invalid InitMessage format:", err)
		conn.Close()
		return
	}

	var userID string
	var errJWT error

	if userID, errJWT = security.DecodeJWT(init.Token, g.jwtSecretKey); errJWT != nil {
		log.Println(errJWT)
		conn.Close()
		return
	}

	client := NewClient(userID, init.RoomID, conn)

	g.mutex.Lock()
	if _, ok := g.rooms[init.RoomID]; !ok {
		g.rooms[init.RoomID] = make(map[string]Client)
	}
	g.rooms[init.RoomID][userID] = client
	g.mutex.Unlock()

	log.Printf("Client %s joined room %s", userID, init.RoomID)

	go g.readPump(client)
	go g.writePump(client)
}

func (g *gatewayImpl) readPump(client Client) {
	defer g.disconnect(client)

	for {
		_, msg, err := client.Read()
		if err != nil {
			log.Printf("Read error [%s]: %v", client.ID(), err)
			break
		}

		to := extractTargetClientID(msg)

		g.mutex.Lock()
		targetRoom, ok := g.rooms[client.Room()]
		targetClient, found := targetRoom[to]
		g.mutex.Unlock()

		if ok && found {
			targetClient.Send(msg)
		} else {
			log.Printf("Target client %s not found in room %s", to, client.Room())
		}
	}
}

func (g *gatewayImpl) writePump(client Client) {
	defer close(client.Buffer())

	for msg := range client.Buffer() {
		if err := client.WriteTextMessage(msg); err != nil {
			log.Printf("Write error [%s]: %v", client.ID(), err)
			break
		}
	}
}

func (g *gatewayImpl) Broadcast(to string, message []byte) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for _, clients := range g.rooms {
		if client, ok := clients[to]; ok {
			client.Send(message)
			return nil
		}
	}
	return fmt.Errorf("client %s not found", to)
}

func (g *gatewayImpl) disconnect(client Client) {
	log.Printf("Disconnecting client %s from room %s", client.ID(), client.Room())

	client.Close()

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if roomClients, ok := g.rooms[client.Room()]; ok {
		delete(roomClients, client.ID())
		if len(roomClients) == 0 {
			delete(g.rooms, client.Room())
			log.Printf("Room %s removed", client.Room())
		}
	}
}

func extractTargetClientID(msg []byte) string {
	type message struct {
		To string `json:"to"`
	}

	var m message
	_ = json.Unmarshal(msg, &m)
	return strings.TrimSpace(m.To)
}
