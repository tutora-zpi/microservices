package ws

import (
	"context"
	"log"
	"sync"
	"ws-gateway/internal/app/interfaces"
	wsevent "ws-gateway/internal/domain/ws_event"

	"github.com/gorilla/websocket"
)

type clientImpl struct {
	id     string
	conn   interfaces.Connection
	send   chan []byte
	done   chan struct{}
	cancel context.CancelFunc
	closed sync.Once
	hub    interfaces.HubManager
}

func NewClient(id string, conn *websocket.Conn, h interfaces.HubManager) interfaces.Client {
	c := &clientImpl{
		id:   id,
		conn: NewConnection(conn),
		send: make(chan []byte, 256),
		done: make(chan struct{}),
		hub:  h,
	}

	go c.runSendLoop()

	return c
}

func (c *clientImpl) ID() string { return c.id }

func (c *clientImpl) GetConnection() interfaces.Connection { return c.conn }

func (c *clientImpl) Close() {
	c.closed.Do(func() {
		close(c.done)
		if c.cancel != nil {
			c.cancel()
		}
		c.hub.RemoveGlobalMember(c)
		c.conn.Close()
	})
}

func (c *clientImpl) runSendLoop() {
	defer func() { recover() }()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.Close()
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *clientImpl) Listen(ctx context.Context, handler func(context.Context, string, []byte, interfaces.Client) error) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	defer c.Close()

	for {
		select {
		case <-ctx.Done():
		case <-c.done:
			return
		default:
			mt, msg, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			if mt != websocket.TextMessage {
				continue
			}

			wrapper, err := wsevent.DecodeSocketEventWrapper(msg)
			if err != nil {
				continue
			}

			if err := handler(ctx, wrapper.Name, wrapper.Payload, c); err != nil {
				log.Printf("handler err %s: %v", c.id, err)
			}
		}
	}
}
