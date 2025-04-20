package server

import (
	"fmt"
	"image"
	"io"

	"github.com/adamroach/webrd/pkg/capture"
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/codec/openh264"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
)

type VideoReader struct {
	capturer capture.VideoCapturer
	image    image.Image
}

func (r *VideoReader) waitForImage() {
	// Wait for an image to be available
	r.image = <-r.capturer.FrameChannel()
}

func (r *VideoReader) Read() (img image.Image, release func(), err error) {
	release = func() {}
	if r.image == nil {
		r.waitForImage()
	}
	if r.image == nil {
		err = io.EOF
	}
	img = r.image
	r.image = nil
	return
}

type VideoEncoder struct {
	reader    *VideoReader
	encoder   codec.ReadCloser
	bitrate   int
	framerate int
	width     int
	height    int
}

func NewVideoEncoder(capturer capture.VideoCapturer, bitrate int, framerate int) (*VideoEncoder, error) {
	r := &VideoEncoder{
		reader:    &VideoReader{capturer: capturer},
		bitrate:   bitrate,
		framerate: framerate,
	}
	return r, nil
}

func (e *VideoEncoder) Read() (b []byte, release func(), err error) {
	release = func() {}
	e.reader.waitForImage()
	if e.reader.image == nil {
		err = io.EOF
		return
	}
	if e.encoder == nil || e.width != e.reader.image.Bounds().Dx() || e.height != e.reader.image.Bounds().Dy() {
		e.Close()
		e.encoder = nil
		e.width = e.reader.image.Bounds().Dx()
		e.height = e.reader.image.Bounds().Dy()
		fmt.Printf("Initializing H.264 encoder: %d x %d @ %vfps\n", e.width, e.height, e.framerate)
		params, _ := openh264.NewParams()
		params.BitRate = e.bitrate
		params.EnableFrameSkip = false
		// Suppress automatic keyframe generation
		params.IntraPeriod = 0

		mediaProperties := prop.Media{
			Video: prop.Video{
				Width:       e.width,
				Height:      e.height,
				FrameRate:   float32(e.framerate),
				FrameFormat: frame.FormatI420,
			},
		}

		e.encoder, err = params.BuildVideoEncoder(e.reader, mediaProperties)
		if err != nil {
			return
		}
	}
	return e.encoder.Read()
}

func (e *VideoEncoder) Close() error {
	if e.encoder == nil {
		return nil
	}
	return e.encoder.Close()
}

func (e *VideoEncoder) Controller() codec.EncoderController {
	if e.encoder == nil {
		return nil
	}
	return e.encoder.Controller()
}
