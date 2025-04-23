package server

import (
	"fmt"
	"log"

	"github.com/adamroach/webrd/pkg/config"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"
)

type WebRTCConnection struct {
	pc          *webrtc.PeerConnection
	videoSender Sender
	audioSender Sender
	iceServers  []webrtc.ICEServer
}

func NewWebRTCConnection(opts ...func(*WebRTCConnection) error) (*WebRTCConnection, error) {
	c := &WebRTCConnection{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	me := &webrtc.MediaEngine{}
	if c.videoSender != nil {
		err := c.videoSender.RegisterCodecs(me)
		if err != nil {
			return nil, fmt.Errorf("error registering video codecs: %v", err)
		}
	}
	if c.audioSender != nil {
		err := c.audioSender.RegisterCodecs(me)
		if err != nil {
			return nil, fmt.Errorf("error registering audio codecs: %v", err)
		}
	}

	ir := &interceptor.Registry{}
	err := webrtc.RegisterDefaultInterceptors(me, ir)
	if err != nil {
		return nil, fmt.Errorf("error registering default interceptors: %v", err)
	}

	se := webrtc.SettingEngine{}
	pcConfig := webrtc.Configuration{}
	if len(c.iceServers) > 0 {
		pcConfig.ICEServers = c.iceServers
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(me), webrtc.WithInterceptorRegistry(ir), webrtc.WithSettingEngine(se))
	c.pc, err = api.NewPeerConnection(pcConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating peer connection: %v", err)
	}

	c.pc.OnConnectionStateChange(c.HandleConnectionStateChange)

	return c, nil
}

func (c *WebRTCConnection) HandleConnectionStateChange(state webrtc.PeerConnectionState) {
	switch state {
	case webrtc.PeerConnectionStateConnected:
		err := c.start()
		if err != nil {
			log.Printf("error starting connection: %v", err)
			// TODO handle failure
		}
	case webrtc.PeerConnectionStateDisconnected:
		// TODO Handle disconnected state
	case webrtc.PeerConnectionStateFailed:
		// TODO Handle failed state
	case webrtc.PeerConnectionStateClosed:
		// TODO Handle closed state
	}
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

func WithICEServers(iceServers []config.IceServer) func(c *WebRTCConnection) error {
	return func(c *WebRTCConnection) error {
		for _, server := range iceServers {
			iceServer := webrtc.ICEServer{
				URLs: server.Urls,
			}
			if server.Username != nil {
				iceServer.Username = *server.Username
			}
			if server.Credential != nil {
				iceServer.Credential = *server.Credential
				iceServer.CredentialType = webrtc.ICECredentialTypePassword
			}
			c.iceServers = append(c.iceServers, iceServer)
		}
		return nil
	}
}

func (c *WebRTCConnection) GetOffer() (string, error) {
	if c.audioSender != nil {
		log.Printf("Adding video track")
		err := c.audioSender.AddTrack(c.pc)
		if err != nil {
			return "", fmt.Errorf("error adding audio track: %v", err)
		}
	}
	if c.videoSender != nil {
		log.Printf("Adding video track")
		err := c.videoSender.AddTrack(c.pc)
		if err != nil {
			return "", fmt.Errorf("error adding video track: %v", err)
		}
	}
	offer, err := c.pc.CreateOffer(nil)
	if err != nil {
		return "", fmt.Errorf("error creating offer: %v", err)
	}
	err = c.pc.SetLocalDescription(offer)
	if err != nil {
		return "", fmt.Errorf("error setting local description: %v", err)
	}
	<-webrtc.GatheringCompletePromise(c.pc)
	return c.pc.LocalDescription().SDP, nil
}

func (c *WebRTCConnection) SetAnswer(answer string) error {
	parsedAnswer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answer,
	}
	err := c.pc.SetRemoteDescription(parsedAnswer)
	if err != nil {
		return fmt.Errorf("error setting remote description: %v", err)
	}
	return nil
}

func (c *WebRTCConnection) AddICECandidate(candidate Candidate) error {
	mlineIndex := uint16(candidate.SdpMLineIndex)
	parsedCandidate := webrtc.ICECandidateInit{
		Candidate:        candidate.Candidate,
		SDPMLineIndex:    &mlineIndex,
		SDPMid:           &candidate.SdpMid,
		UsernameFragment: &candidate.UsernameFragment,
	}
	err := c.pc.AddICECandidate(parsedCandidate)
	if err != nil {
		return fmt.Errorf("error adding ICE candidate: %v", err)
	}
	return nil
}

func (c *WebRTCConnection) start() error {
	if c.audioSender != nil {
		err := c.audioSender.Start()
		if err != nil {
			return fmt.Errorf("error starting audio sender: %v", err)
		}
	}
	if c.videoSender != nil {
		err := c.videoSender.Start()
		if err != nil {
			return fmt.Errorf("error starting video sender: %v", err)
		}
	}
	return nil
}

func (c *WebRTCConnection) Close() error {
	if c.audioSender != nil {
		err := c.audioSender.Close()
		if err != nil {
			return fmt.Errorf("error closing audio sender: %v", err)
		}
	}
	if c.videoSender != nil {
		err := c.videoSender.Close()
		if err != nil {
			return fmt.Errorf("error closing video sender: %v", err)
		}
	}
	return c.pc.Close()
}
