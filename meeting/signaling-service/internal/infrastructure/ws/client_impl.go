package ws

import (
	"context"
	"log"
	"signaling-service/internal/app/interfaces"
	wsevent "signaling-service/internal/domain/ws_event"

	"github.com/gorilla/websocket"
)

type clientImpl struct {
	id     string
	conn   interfaces.Connection
	cancel context.CancelFunc
}

// GetConnection implements interfaces.Client.
func (c *clientImpl) GetConnection() interfaces.Connection {
	return c.conn
}

type connImpl struct {
	conn *websocket.Conn
}

// WriteMessage implements models.Connection.
func (c *connImpl) WriteMessage(messageType int, payload []byte) error {
	return c.conn.WriteMessage(messageType, payload)
}

// Close implements models.Connection.
func (c *connImpl) Close() {
	c.conn.Close()
}

// ReadMessage implements models.Connection.
func (c *connImpl) ReadMessage() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

func NewConnection(conn *websocket.Conn) interfaces.Connection {
	return &connImpl{
		conn: conn,
	}
}

func NewClient(id string, conn *websocket.Conn) interfaces.Client {
	return &clientImpl{
		id:   id,
		conn: NewConnection(conn),
	}
}

// ID implements Client.
func (c *clientImpl) ID() string {
	return c.id
}

// Listen implements Client.
func (c *clientImpl) Listen(ctx context.Context, handler func(ctx context.Context, eventType string, msg []byte, client interfaces.Client) error) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	defer cancel()
	defer c.GetConnection().Close()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context finished")
			return
		default:
			messageType, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("Read error from %s: %v", c.ID(), err)
				return
			}

			switch messageType {
			case websocket.TextMessage:

				wrapper, err := wsevent.DecodeSocketEventWrapper(msg)

				if err != nil {
					log.Printf("Something went wrong: %v", err)
					continue
				}

				err = handler(ctx, wrapper.Name, wrapper.Payload, c)
				if err != nil {
					log.Println(err)
				}

			case websocket.CloseMessage:
				log.Printf("Client %s closed connection", c.ID())
				return
			case websocket.PingMessage:
				return
			default:
				log.Println("Unsupported messageType")
				continue
			}
		}
	}
}
