package peer

import (
	"context"
	"log"
	"recorder-service/internal/domain/recorder"
	wsevent "recorder-service/internal/domain/ws_event"
	"recorder-service/internal/domain/ws_event/rtc"
	"recorder-service/internal/infrastructure/webrtc/writer"

	"github.com/pion/webrtc/v3"
)

type Peer interface {
	Init() error
	CreateOffer() error
	CreateAnswer() error
	SetRemoteDescription(desc webrtc.SessionDescription) error
	AddICECandidate(candidate webrtc.ICECandidateInit) error
	HandleTrack(track *webrtc.TrackRemote)
	StopReceivingTracks() *recorder.RecordingInfo
	Close() error
	GetBotID() string
	GetRemoteUserID() string
}

type PeerImpl struct {
	BotID      string
	RemoteUser string
	RoomID     string
	peerConn   *webrtc.PeerConnection
	recorder   recorder.Recorder
	sendWS     func(msg []byte) error
}

// StopReceivingTracks implements Peer.
func (p *PeerImpl) StopReceivingTracks() *recorder.RecordingInfo {
	return p.recorder.StopRecording(p.RemoteUser)
}

func NewPeer(roomID, botID, remoteUser string, rec recorder.Recorder, sendWS func(msg []byte) error) Peer {
	return &PeerImpl{
		RoomID:     roomID,
		BotID:      botID,
		RemoteUser: remoteUser,
		recorder:   rec,
		sendWS:     sendWS,
	}
}

func (p *PeerImpl) Init() error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		evt := &rtc.IceCandidateWSEvent{
			Candidate: c.ToJSON(),
			RoomID:    p.RoomID,
			From:      p.BotID,
			To:        p.RemoteUser,
		}

		msg, err := wsevent.EncodeSocketEventWrapper(evt)
		if err != nil {
			log.Printf("Failed to encode event: %v", err)
			return
		}

		if err := p.sendWS(msg); err != nil {
			log.Printf("Failed to sent message: %v", err)
		}
	})

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Received track from %s: %s", p.RemoteUser, track.Kind())
		p.HandleTrack(track)
	})

	p.peerConn = pc
	return nil
}

func (p *PeerImpl) CreateOffer() error {
	if p.peerConn == nil {
		return webrtc.ErrConnectionClosed
	}

	offer, err := p.peerConn.CreateOffer(nil)
	if err != nil {
		return err
	}

	if err := p.peerConn.SetLocalDescription(offer); err != nil {
		return err
	}

	log.Println("Offer created and set as local description")

	evt := &rtc.OfferWSEvent{
		Offer: offer,
		From:  p.BotID,
		To:    p.RemoteUser,
	}

	msg, err := wsevent.EncodeSocketEventWrapper(evt)
	if err != nil {
		log.Printf("Failed to encode Offer event: %v", err)
		return err
	}
	if err := p.sendWS(msg); err != nil {
		log.Printf("Failed to send Offer event: %v", err)
		return err
	}

	return nil
}

func (p *PeerImpl) CreateAnswer() error {
	if p.peerConn == nil {
		return webrtc.ErrConnectionClosed
	}

	answer, err := p.peerConn.CreateAnswer(nil)
	if err != nil {
		return err
	}

	if err := p.peerConn.SetLocalDescription(answer); err != nil {
		return err
	}

	log.Println("Answer created and set as local description")

	evt := &rtc.AnswerWSEvent{
		Answer: answer,
		From:   p.BotID,
		To:     p.RemoteUser,
		RoomID: p.RoomID,
	}

	msg, err := wsevent.EncodeSocketEventWrapper(evt)
	if err != nil {
		log.Printf("Failed to encode Answer event: %v", err)
		return err
	}

	if p.sendWS != nil {
		if err := p.sendWS(msg); err != nil {
			log.Printf("Failed to send Answer event: %v", err)
			return err
		}
	}

	return nil
}

func (p *PeerImpl) SetRemoteDescription(desc webrtc.SessionDescription) error {
	return p.peerConn.SetRemoteDescription(desc)
}

func (p *PeerImpl) AddICECandidate(candidate webrtc.ICECandidateInit) error {
	return p.peerConn.AddICECandidate(candidate)
}

func (p *PeerImpl) HandleTrack(track *webrtc.TrackRemote) {
	ctx := context.Background()
	p.recorder.StartRecording(ctx, p.RoomID, p.RemoteUser, track, writer.NewLocalWriter)
}

func (p *PeerImpl) Close() error {
	if p.peerConn != nil {
		log.Printf("Closing peer for %s", p.RemoteUser)
		return p.peerConn.Close()
	}
	return nil
}

func (p *PeerImpl) GetBotID() string {
	return p.BotID
}

func (p *PeerImpl) GetRemoteUserID() string {
	return p.RemoteUser
}
