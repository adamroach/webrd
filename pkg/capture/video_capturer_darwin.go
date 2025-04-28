package capture

import "github.com/adamroach/webrd/pkg/capture/darwin"

func NewVideoCapturer(framerate int) (VideoCapturer, error) {
	return darwin.NewVideoCapturer(framerate)
	//return screenshot.NewVideoCapturer(framerate)
}
