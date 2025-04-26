package server

import (
	"fmt"
	"io"
	"log"

	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/hid"
	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	Server           *Server
	WebRTCConnection *WebRTCConnection
	MessageChannel   MessageChannel
	VideoCapturer    capture.VideoCapturer
	AudioCapturer    capture.AudioCapturer
	Keyboard         hid.Keyboard
	Mouse            hid.Mouse
}

func (s *Session) Start() error {
	offer, err := s.WebRTCConnection.GetOffer()
	if err != nil {
		log.Printf("could not get offer: %v", err)
		return err
	}
	err = s.MessageChannel.Send(OfferMessage{Type: TypeOffer, SDP: offer})
	if err != nil {
		log.Printf("could not send offer: %v", err)
	}
	if s.VideoCapturer != nil {
		if err := s.VideoCapturer.Start(); err != nil {
			log.Printf("could not start video capturer: %v", err)
			return err
		}
	}
	if s.AudioCapturer != nil {
		if err := s.AudioCapturer.Start(); err != nil {
			log.Printf("could not start audio capturer: %v", err)
			return err
		}
	}
	go s.handleMessages()
	return nil
}

func (s *Session) Close() error {
	err := s.MessageChannel.Close()
	if err != nil {
		return err
	}
	if s.VideoCapturer != nil {
		if err := s.VideoCapturer.Stop(); err != nil {
			return fmt.Errorf("could not stop video capturer: %v", err)
		}
	}
	if s.AudioCapturer != nil {
		if err := s.AudioCapturer.Stop(); err != nil {
			return fmt.Errorf("could not stop audio capturer: %v", err)
		}
	}

	s.Server.removeSession(s)
	return nil
}

func (s *Session) handleMessages() {
	for {
		message, err := s.MessageChannel.Receive()
		if err != nil {
			if err == io.EOF {
				log.Printf("session closed: %v\n", err)
				return
			}
			log.Printf("could not receive message: %v\n", err)
			continue
		}

		switch message := message.(type) {
		case *AnswerMessage:
			err = s.WebRTCConnection.SetAnswer(message.SDP)
			if err != nil {
				log.Printf("could not set answer: %v\n", err)
			}
		case *IceCandidateMessage:
			err = s.WebRTCConnection.AddICECandidate(message.Candidate)
			if err != nil {
				log.Printf("could not add ICE candidate: %v\n", err)
			}
		case *KeyboardMessage:
			if s.Keyboard != nil {
				err = s.Keyboard.Key(message.Event)
				if err != nil {
					log.Printf("could not send keyboard event: %v\n", err)
				}
			} else {
				log.Printf("keyboard not available\n")
			}
		case *MouseButtonMessage:
			if s.Mouse != nil {
				err = s.Mouse.Button(message.Button, message.X, message.Y, message.Down)
				if err != nil {
					log.Printf("could not send mouse button event: %v\n", err)
				}
			} else {
				log.Printf("mouse not available\n")
			}
		case *MouseWheelMessage:
			if s.Mouse != nil {
				err = s.Mouse.Wheel(message.DeltaX, message.DeltaY, message.DeltaZ)
				if err != nil {
					log.Printf("could not send mouse wheel event: %v\n", err)
				}
			}
		case *MouseMoveMessage:
			if s.Mouse != nil {
				err = s.Mouse.Move(message.X, message.Y)
				if err != nil {
					log.Printf("could not send mouse move event: %v\n", err)
				}
			}
			// we don't log the "else" clause here because it would be too noisy

		default:
			log.Printf("unexpected message type: %+v\n", message)
		}
	}
}
