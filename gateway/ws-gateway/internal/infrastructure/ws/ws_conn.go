package ws

import (
	"sync"
	"ws-gateway/internal/app/interfaces"

	"github.com/gorilla/websocket"
)

type connImpl struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *connImpl) WriteMessage(mt int, p []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(mt, p)
}

func (c *connImpl) ReadMessage() (int, []byte, error) { return c.conn.ReadMessage() }

func (c *connImpl) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.Close()
}

func NewConnection(c *websocket.Conn) interfaces.Connection { return &connImpl{conn: c} }
