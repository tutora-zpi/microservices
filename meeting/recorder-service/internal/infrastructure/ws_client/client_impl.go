package ws

import (
	"context"
	"fmt"
	"log"
	"recorder-service/internal/domain/client"
	wsevent "recorder-service/internal/domain/ws_event"
	"recorder-service/internal/infrastructure/bus"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type socketClientImpl struct {
	conn      *websocket.Conn
	dispacher bus.Dispachable

	peer *webrtc.PeerConnection

	url string
	cfg webrtc.Configuration
}

// Close implements client.Client.
func (s *socketClientImpl) Close() error {
	err := s.peer.Close()
	if err != nil {
		log.Printf("Failed to close peer: %v", err)
	}

	err = s.conn.Close()
	if err != nil {
		log.Printf("Failed to close websocket connection: %v", err)
	}

	return nil
}

// Connect implements client.Client.
func (s *socketClientImpl) Connect(ctx context.Context) error {
	peer, err := webrtc.NewPeerConnection(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to create new peer connection: %w", err)
	}

	s.peer = peer

	var dialer websocket.Dialer
	conn, _, err := dialer.Dial(s.url, nil)

	if err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}

	s.conn = conn

	go s.Listen(ctx)

	return nil
}

// OnTrack implements client.Client.
func (s *socketClientImpl) OnTrack(callback func(*webrtc.TrackRemote)) {
	s.peer.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		callback(track)
	})
}

// Send implements client.Client.
func (s *socketClientImpl) Send(msg []byte) error {
	err := s.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Printf("Failed to write message: %v", err)
	}
	return nil
}

// ValidateMessageType implements Client.
func (s *socketClientImpl) ValidateMessageType(msgType int) error {
	if msgType != websocket.TextMessage {
		return fmt.Errorf("unsupported message type: %d", msgType)
	}
	return nil
}

// Listen implements SocketClient.
func (s *socketClientImpl) Listen(ctx context.Context) {
	for {
		msgType, body, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message, skipping")
			continue
		}

		err = s.ValidateMessageType(msgType)
		if err != nil {
			log.Printf("Invalid message type: %v", err)
			continue
		}

		wrapper, err := wsevent.DecodeSocketEventWrapper(body)
		if err != nil {
			log.Printf("Failed to decode body into valid wrapper: %v", err)
			continue
		}

		if err := s.dispacher.HandleEvent(ctx, wrapper.Name, wrapper.ToBytes()); err != nil {
			log.Printf("An error occurred during event handling: %v", err)
			continue
		}
	}
}

func NewSocketClient(url string, webrtcConfig *webrtc.Configuration) client.Client {
	cfg := webrtcConfig
	if cfg == nil {
		cfg = &webrtc.Configuration{}
	}

	return &socketClientImpl{url: url, cfg: *cfg}
}
