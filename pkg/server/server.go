package server

import (
	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/hid"
)

type Server struct {
	MakeVideoCapturer func() (capture.VideoCapturer, error)
	MakeKeyboard      func() (hid.Keyboard, error)
	MakeMouse         func() (hid.Mouse, error)
}
