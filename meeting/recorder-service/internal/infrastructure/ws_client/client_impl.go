package wsclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"recorder-service/internal/domain/client"
	wsevent "recorder-service/internal/domain/ws_event"
	"recorder-service/internal/infrastructure/bus"
	"recorder-service/internal/infrastructure/security"
	"strings"

	"sync"

	"github.com/gorilla/websocket"
)

type clientImpl struct {
	dispatcher bus.Dispachable
	conn       *websocket.Conn
	botID      string
	msg        chan []byte
	url        string

	mu       sync.Mutex
	isClosed bool
}

func NewWSClient(url string, dispatcher bus.Dispachable) client.Client {
	return &clientImpl{
		url:        url,
		dispatcher: dispatcher,
		msg:        make(chan []byte, 256),
	}
}

func (c *clientImpl) SetBotID(botID string) {
	c.botID = botID
}

func (c *clientImpl) Connect(ctx context.Context) error {
	var err error

	token, err := security.FetchToken(ctx, c.botID)
	if err != nil {
		return err
	}

	header := http.Header{}
	bearer := strings.Join([]string{"Bearer", token.AccessToken}, " ")
	header.Set("Authorization", bearer)

	c.conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("%s", c.url), header)
	if err != nil {
		return err
	}

	go c.listen(ctx)
	go c.sending(ctx)

	return nil
}

func (c *clientImpl) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return nil
	}

	c.isClosed = true

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *clientImpl) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn != nil && c.conn.RemoteAddr() != nil && !c.isClosed
}

func (c *clientImpl) Send(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isClosed {
		return fmt.Errorf("client is closed")
	}

	select {
	case c.msg <- msg:
		return nil
	default:
		return fmt.Errorf("send buffer full")
	}
}

func (c *clientImpl) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			decoded, err := wsevent.DecodeSocketEventWrapper(msg)
			if err != nil {
				log.Println("Failed to decode event:", err)
				continue
			}

			if err := c.dispatcher.HandleEvent(ctx, decoded.Name, decoded.Payload); err != nil {
				log.Printf("Failed to handle event: %v", err)
			}
		}
	}
}

func (c *clientImpl) sending(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.msg:
			c.mu.Lock()
			if c.isClosed {
				c.mu.Unlock()
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			c.mu.Unlock()
			if err != nil {
				log.Println("Write error:", err)
				return
			}
		}
	}
}
