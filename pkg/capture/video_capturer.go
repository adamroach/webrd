package capture

import "image"

type VideoCapturer interface {
	Start() error
	Stop() error
	GetBounds() image.Rectangle
	FrameChannel() <-chan image.Image // TODO -- include timestamp information
}
