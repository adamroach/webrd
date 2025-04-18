package server

import "github.com/pion/webrtc/v4"

type Sender interface {
	RegisterCodecs(me *webrtc.MediaEngine) error
	AddTrack(pc *webrtc.PeerConnection) error
	Start() error
	Close() error
}
