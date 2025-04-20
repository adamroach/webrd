package main

import (
	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/hid"
	"github.com/adamroach/webrd/pkg/server"
)

func main() {
	server := server.Server{
		MakeVideoCapturer: func() (capture.VideoCapturer, error) {
			return capture.NewVideoCapturer()
		},
		MakeAudioCapturer: nil,
		MakeKeyboard: func() (hid.Keyboard, error) {
			return hid.NewKeyboard()
		},
		MakeMouse: nil,
	}
	panic(server.Run())
}
