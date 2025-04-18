package server

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
)

type VideoSender struct {
	encoder            codec.ReadCloser
	track              *webrtc.TrackLocalStaticRTP
	sender             *webrtc.RTPSender
	packetizer         rtp.Packetizer
	codecCapability    webrtc.RTPCodecCapability
	keyFrameController codec.KeyFrameController
}

func NewVideoSender(encoder codec.ReadCloser, keyFrameController codec.KeyFrameController) *VideoSender {
	return &VideoSender{
		encoder:            encoder,
		keyFrameController: keyFrameController,
	}
}

func (s *VideoSender) RegisterCodecs(me *webrtc.MediaEngine) error {
	s.codecCapability = webrtc.RTPCodecCapability{
		MimeType:    webrtc.MimeTypeH264,
		ClockRate:   90000,
		Channels:    0,
		SDPFmtpLine: "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f",
		RTCPFeedback: []webrtc.RTCPFeedback{
			{Type: "goog-remb", Parameter: ""},
			{Type: "ccm", Parameter: "fir"},
			{Type: "nack", Parameter: ""},
		},
	}
	for _, codec := range []webrtc.RTPCodecParameters{
		{
			RTPCodecCapability: s.codecCapability,
			PayloadType:        102,
		},
		{
			RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: "video/rtx", ClockRate: 90000, Channels: 0, SDPFmtpLine: "apt=102", RTCPFeedback: nil},
			PayloadType:        121,
		},
	} {
		if err := me.RegisterCodec(codec, webrtc.RTPCodecTypeVideo); err != nil {
			return err
		}
	}
	return nil
}

func (s *VideoSender) AddTrack(pc *webrtc.PeerConnection) error {
	var err error
	s.track, err = webrtc.NewTrackLocalStaticRTP(s.codecCapability, "video", "screen")
	if err != nil {
		return err
	}
	s.sender, err = pc.AddTrack(s.track)
	if err != nil {
		return err
	}
	return nil
}

func (s *VideoSender) Start() error {
	s.packetizer = rtp.NewPacketizer(
		1400,
		102,
		rand.Uint32(),
		&codecs.H264Payloader{},
		rtp.NewRandomSequencer(),
		s.codecCapability.ClockRate,
	)

	go s.handleRtcp()
	go s.sendMedia()
	return nil
}

func (s *VideoSender) Close() error {
	if s.encoder != nil {
		if err := s.encoder.Close(); err != nil {
			return err
		}
	}
	if s.sender != nil {
		if err := s.sender.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (s *VideoSender) sendMedia() {
	lastFrameTime := time.Now()
	for {
		f, release, err := s.encoder.Read()
		if err != nil {
			log.Printf("Error reading frame: %v", err)
			return
		}
		delta := time.Since(lastFrameTime)
		lastFrameTime = time.Now()
		samples := int(90000 * delta.Seconds())

		rtpPackets := s.packetizer.Packetize(f, uint32(samples))
		for _, pkt := range rtpPackets {
			buffer, err := pkt.Marshal()
			if err != nil {
				log.Printf("Error marshaling packet: %v", err)
				return
			}
			_, err = s.track.Write(buffer)
			if err != nil {
				log.Printf("Error writing packet: %v", err)
				return
			}
		}
		release()
	}
}

func (s *VideoSender) handleRtcp() {
	buf := make([]byte, 2000)
	for {
		n, _, err := s.sender.Read(buf)
		if err != nil {
			log.Printf("Error reading RTCP: %v", err)
			return
		}
		messages, err := rtcp.Unmarshal(buf[:n])
		if err != nil {
			log.Printf("Error unmarshalling RTCP: %v", err)
			continue
		}
		for _, message := range messages {
			switch msg := message.(type) {
			case *rtcp.PictureLossIndication:
				log.Printf("Received PLI: %v", msg)
				if s.keyFrameController != nil {
					s.keyFrameController.ForceKeyFrame()
				}
			case *rtcp.FullIntraRequest:
				log.Printf("Received FIR: %v", msg)
				if s.keyFrameController != nil {
					s.keyFrameController.ForceKeyFrame()
				}
			default:
				// Handle other RTCP messages if needed
			}
		}
	}
}
