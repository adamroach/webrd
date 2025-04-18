package server

import (
	"fmt"

	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/hid"
	"github.com/google/uuid"
)

type Server struct {
	MakeVideoCapturer func() (capture.VideoCapturer, error)
	MakeKeyboard      func() (hid.Keyboard, error)
	MakeMouse         func() (hid.Mouse, error)
	sessions          map[uuid.UUID]*Session // this probably needs a mutex?
}

func (s *Server) NewSession(messageChannel MessageChannel) (*Session, error) {
	videoCapturer, err := s.MakeVideoCapturer()
	if err != nil {
		return nil, fmt.Errorf("could not create video capturer: %v", err)
	}
	if err := videoCapturer.Start(); err != nil {
		return nil, fmt.Errorf("could not start video capturer: %v", err)
	}

	keyboard, err := s.MakeKeyboard()
	if err != nil {
		return nil, fmt.Errorf("could not create keyboard: %v", err)
	}

	mouse, err := s.MakeMouse()
	if err != nil {
		return nil, fmt.Errorf("could not create mouse: %v", err)
	}

	session := &Session{
		ID:               uuid.New(),
		WebRTCConnection: &WebRTCConnection{},
		MessageChannel:   messageChannel,
		VideoCapturer:    videoCapturer,
		Keyboard:         keyboard,
		Mouse:            mouse,
	}

	err = session.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start session: %v", err)
	}

	if s.sessions == nil {
		s.sessions = make(map[uuid.UUID]*Session)
	}
	s.sessions[session.ID] = session

	return session, nil
}
