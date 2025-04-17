package server

import (
	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/hid"
)

type Session struct {
	WebRTCConnection *WebRTCConnection
	MessageChannel   MessageChannel
	VideoCapturer    capture.VideoCapturer
	Keyboard         hid.Keyboard
	Mouse            hid.Mouse
}
