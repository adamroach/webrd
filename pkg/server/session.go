package server

import (
	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/hid"
	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	WebRTCConnection *WebRTCConnection
	MessageChannel   MessageChannel
	VideoCapturer    capture.VideoCapturer
	Keyboard         hid.Keyboard
	Mouse            hid.Mouse
}

func (s *Session) Start() error {
	return nil
}

func (s *Session) Close() error {
	err := s.MessageChannel.Close()
	return err
}
