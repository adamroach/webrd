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

type AudioSender struct {
	encoder         codec.ReadCloser
	track           *webrtc.TrackLocalStaticRTP
	sender          *webrtc.RTPSender
	packetizer      rtp.Packetizer
	codecCapability webrtc.RTPCodecCapability
}

func (s *AudioSender) RegisterCodecs(me *webrtc.MediaEngine) error {
	s.codecCapability = webrtc.RTPCodecCapability{
		MimeType:    webrtc.MimeTypeOpus,
		ClockRate:   48000,
		Channels:    2,
		SDPFmtpLine: "minptime=10;useinbandfec=1",
		RTCPFeedback: []webrtc.RTCPFeedback{
			{Type: "nack", Parameter: ""},
		},
	}
	for _, codec := range []webrtc.RTPCodecParameters{
		{
			RTPCodecCapability: s.codecCapability,
			PayloadType:        111,
		},
	} {
		if err := me.RegisterCodec(codec, webrtc.RTPCodecTypeAudio); err != nil {
			return err
		}
	}
	return nil
}

func (s *AudioSender) AddTrack(pc *webrtc.PeerConnection) error {
	var err error
	s.track, err = webrtc.NewTrackLocalStaticRTP(s.codecCapability, "audio", "speaker")
	if err != nil {
		return err
	}
	s.sender, err = pc.AddTrack(s.track)
	if err != nil {
		return err
	}
	return nil
}

func (s *AudioSender) Start() error {
	s.packetizer = rtp.NewPacketizer(
		1400,
		111,
		rand.Uint32(),
		&codecs.OpusPayloader{},
		rtp.NewRandomSequencer(),
		s.codecCapability.ClockRate,
	)

	go s.handleRtcp()
	go s.sendMedia()
	return nil
}

func (s *AudioSender) Close() error {
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

func (s *AudioSender) sendMedia() {
	lastFrameTime := time.Now()
	for {
		f, release, err := s.encoder.Read()
		if err != nil {
			log.Printf("Error reading audio: %v", err)
			return
		}
		// TODO This logic isn't correct for audio -- the timestamps should be
		// precisely calculated based on the audio frame size and sample rate
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

func (s *AudioSender) handleRtcp() {
	// Even if we're not acting on them, we need to read RTCP packets
	// so that the NACK interceptor can run
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
			default:
				// Handle RTCP messages if needed
				_ = msg
			}
		}
	}
}
