package main

import (
	"github.com/adamroach/webrd/pkg/auth"
	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/config"
	"github.com/adamroach/webrd/pkg/hid"
	"github.com/adamroach/webrd/pkg/server"
)

func main() {
	config := config.NewConfig()

	var authenticator auth.Authenticator

	if config.Auth.UseSystemAuth {
		authenticator = auth.NewSystemAuthenticator(&config.Auth)
	} else {
		authenticator = auth.NewStaticAuthenticator(&config.Auth)
	}

	server := server.Server{
		MakeVideoCapturer: func() (capture.VideoCapturer, error) {
			return capture.NewVideoCapturer(config.Video.Framerate)
		},
		MakeAudioCapturer: nil,
		MakeKeyboard: func() (hid.Keyboard, error) {
			return hid.NewKeyboard()
		},
		MakeMouse: func() (hid.Mouse, error) {
			return hid.NewMouse()
		},
		Authenticator: authenticator,
	}
	panic(server.Run(config))
}
