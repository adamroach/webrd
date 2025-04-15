package capture

import "github.com/adamroach/webrd/pkg/capture/darwin"

func NewVideoCapturer() VideoCapturer {
	return darwin.NewVideoCapturer()
}
