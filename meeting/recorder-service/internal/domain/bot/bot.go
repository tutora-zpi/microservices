package bot

import (
	"fmt"
	"log"
	"maps"
	"recorder-service/internal/domain/client"
	"recorder-service/internal/domain/dto"
	"recorder-service/internal/domain/recorder"
	"recorder-service/internal/infrastructure/webrtc/peer"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type Bot interface {
	ID() string
	Name() string

	DTO() dto.BotDTO

	Client() client.Client

	Peers() map[string]peer.Peer
	AddPeer(remoteUserID string, p peer.Peer) error
	GetPeer(remoteUserID string) (peer.Peer, bool)
	FinishRecording() []recorder.RecordingInfo
}

type bot struct {
	id   string
	name string

	peersMu sync.Mutex
	peers   map[string]peer.Peer

	client client.Client
}

// FinishRecording implements Bot.
func (b *bot) FinishRecording() []recorder.RecordingInfo {
	infos := []recorder.RecordingInfo{}
	for userID, peer := range b.peers {
		log.Printf("Stopping for %s", userID)
		info := peer.StopReceivingTracks()
		if info != nil {
			infos = append(infos, *info)
		}
		peer.Close()
	}

	b.client.Close()

	return infos
}

func (b *bot) Peers() map[string]peer.Peer {
	b.peersMu.Lock()
	defer b.peersMu.Unlock()

	copy := make(map[string]peer.Peer, len(b.peers))
	maps.Copy(copy, b.peers)
	return copy
}

func (b *bot) AddPeer(remoteUserID string, p peer.Peer) error {
	if err := p.Init(); err != nil {
		return fmt.Errorf("failed to init peer for user %s: %w", remoteUserID, err)
	}

	b.peersMu.Lock()
	defer b.peersMu.Unlock()

	if b.peers == nil {
		b.peers = make(map[string]peer.Peer)
	}

	b.peers[remoteUserID] = p
	return nil
}

func (b *bot) GetPeer(remoteUserID string) (peer.Peer, bool) {
	b.peersMu.Lock()
	defer b.peersMu.Unlock()

	p, ok := b.peers[remoteUserID]
	return p, ok
}

func (b *bot) DTO() dto.BotDTO {
	return dto.BotDTO{
		ID:   b.id,
		Name: b.name,
	}
}

func (b *bot) Client() client.Client {
	return b.client
}

func (b *bot) ID() string {
	return b.id
}

func (b *bot) Name() string {
	return b.name
}

func NewBot(client client.Client) Bot {
	botID := uuid.NewString()

	client.SetBotID(botID)

	return &bot{
		id:     botID,
		name:   gofakeit.FirstName(),
		client: client,
		peers:  make(map[string]peer.Peer),
	}
}
