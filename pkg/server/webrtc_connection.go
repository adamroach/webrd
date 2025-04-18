package server

import (
	"github.com/pion/webrtc/v4"
)

type WebRTCConnection struct {
	pc          *webrtc.PeerConnection
	videoSender *VideoSender
	audioSender *AudioSender
	turnServers []string
}

func NewWebRTCConnection(opts ...func(*WebRTCConnection) error) (*WebRTCConnection, error) {
	c := &WebRTCConnection{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func WithVideoSender(sender *VideoSender) func(c *WebRTCConnection) error {
	return func(c *WebRTCConnection) error {
		c.videoSender = sender
		return nil
	}
}

func WithAudioSender(sender *AudioSender) func(c *WebRTCConnection) error {
	return func(c *WebRTCConnection) error {
		c.audioSender = sender
		return nil
	}
}

func WithTURNServers(turnServers []string) func(c *WebRTCConnection) error {
	return func(c *WebRTCConnection) error {
		c.turnServers = turnServers
		return nil
	}
}
