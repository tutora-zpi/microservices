package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"recorder-service/internal/app/interfaces"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/infrastructure/security"
	"sync"

	"github.com/gorilla/websocket"
)

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

	dispatcher interfaces.Dispatcher
}

func NewGateway(secret string, dispatcher interfaces.Dispatcher) Gateway {
	return &gatewayImpl{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		rooms:        make(map[string]map[string]Client),
		jwtSecretKey: []byte(secret),
		dispatcher:   dispatcher,
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

		action, err := getActionFromMsg(msg)
		if err != nil {
			log.Println("Failed to get action")
		}

		if err := g.dispatcher.Dispatch(action, msg); err != nil {
			log.Printf("Dispatch error [%s]: %v", action, err)
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

func getActionFromMsg(msg []byte) (string, error) {

	var m dto.WsMessageDTO
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return "", err
	}
	return m.Action, nil
}
