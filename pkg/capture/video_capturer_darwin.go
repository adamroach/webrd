package capture

import "github.com/adamroach/webrd/pkg/capture/darwin"

func NewVideoCapturer() (VideoCapturer, error) {
	return darwin.NewVideoCapturer()
}
