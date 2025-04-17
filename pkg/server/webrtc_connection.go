package server

import (
	"github.com/pion/webrtc/v4"
)

type WebRTCConnection struct {
	pc          *webrtc.PeerConnection
	videoSender *VideoSender
	audioSender *AudioSender
}
