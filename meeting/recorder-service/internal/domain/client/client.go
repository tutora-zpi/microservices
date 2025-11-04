package client

import (
	"context"

	"github.com/pion/webrtc/v3"
)

type Client interface {
	Listen(ctx context.Context)
	ValidateMessageType(msgType int) error
	OnTrack(func(*webrtc.TrackRemote))
	Send(msg []byte) error
	Connect(ctx context.Context) error
	Close() error

	SetRemoteDescription(desc webrtc.SessionDescription) error
	AddIceCandidate(candidate webrtc.ICECandidateInit) error
}
