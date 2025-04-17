package server

import (
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type AudioSender struct {
	encoder    codec.ReadCloser
	track      *webrtc.TrackLocalStaticRTP
	sender     *webrtc.RTPSender
	packetizer rtp.Packetizer
}
