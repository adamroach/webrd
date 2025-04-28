package screenshot

import (
	"image"
	"log"
	"time"

	"github.com/adamroach/webrd/pkg/imageconvert"
	"github.com/kbinani/screenshot"
)

type VideoCapturer struct {
	frames       chan (image.Image)
	stop         chan (struct{})
	screenNumber int
	framerate    int
}

func NewVideoCapturer(framerate int) (*VideoCapturer, error) {
	c := &VideoCapturer{
		frames:    make(chan image.Image, 4),
		framerate: framerate,
	}
	return c, nil
}

func (c *VideoCapturer) Start() error {
	duration := time.Duration(float64(1*time.Second) / float64(c.framerate))
	lastFrame := time.Now()
	bounds := screenshot.GetDisplayBounds(c.screenNumber)
	yuvImage := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio420)
	go func() {
		for {
			timeToNextFrame := max(time.Until(lastFrame.Add(duration)), time.Nanosecond)
			lastFrame = time.Now()
			timer := time.NewTimer(timeToNextFrame)
			select {
			case <-c.stop:
				log.Printf("Stopping video capture loop")
				return
			case <-timer.C:
				if len(c.frames) == cap(c.frames) {
					// If the queue is full, we're producing faster than
					// frames can be consumed, so we drop this frame
					continue
				}
				rgbImage, err := screenshot.CaptureRect(bounds)
				if err != nil {
					log.Printf("Error grabbing screenshot; exiting video capture loop: %v", err)
				}
				imageconvert.ToYCbCr(yuvImage, rgbImage)
				c.frames <- yuvImage
			}
		}
	}()
	return nil
}

func (c *VideoCapturer) Stop() error {
	close(c.stop)
	return nil
}

func (c *VideoCapturer) FrameChannel() <-chan image.Image {
	return c.frames
}
