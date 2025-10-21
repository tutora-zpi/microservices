package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client interface {

	// UUID from users-service db
	ID() string

	// Getter for room (meetingID)
	Room() string

	// Sends data to ws client
	Send(data []byte)

	// Closes connection with ws
	Close() error

	// Getter for message
	// Return: message type, data, or error
	Read() (int, []byte, error)

	// Closes message burffer
	Buffer() chan []byte

	WriteTextMessage(msg []byte) error
}

func NewClient(id, room string, conn *websocket.Conn) Client {
	return &clientImpl{
		id:     id,
		room:   room,
		conn:   conn,
		buffer: make(chan []byte, 256),
	}
}

type clientImpl struct {
	// Users ID (UUID) from users-service db
	id string

	conn *websocket.Conn

	// Buffer for outgoing messages (bytes channel)
	buffer chan []byte

	// aka meetingID
	room string
}

// Room implements Client.
func (c *clientImpl) Room() string {
	return c.room
}

// Getter for clients ID
func (c *clientImpl) ID() string {
	return c.id
}

// Sends data to the websocket client
func (c *clientImpl) Send(msg []byte) {
	select {
	case c.buffer <- msg:
	default:
		log.Println("Send channel full, dropping message for client:", c.id)
	}
}

// Closes the websocket connection
func (c *clientImpl) Close() error {
	return c.conn.Close()
}

func (c *clientImpl) Read() (int, []byte, error) {
	return c.conn.ReadMessage()
}

func (c *clientImpl) Buffer() chan []byte {
	return c.buffer
}

func (c *clientImpl) WriteTextMessage(msg []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}
