package capture

type AudioCapturer interface {
	Start() error
	Stop() error
	FrameChannel() <-chan []byte // TODO -- include timestamp information
}
