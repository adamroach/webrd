package capture

import "image"

type VideoCapturer interface {
	Start() error
	Stop() error
	FrameChannel() <-chan image.Image // TODO -- include timestamp information
}
