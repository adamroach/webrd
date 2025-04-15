//go:build cgo && darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework AVFoundation -framework Foundation -framework CoreMedia -framework CoreVideo

#import "video_capturer.h"

static void* newVideoCapturer(){
    void *capturer = [[VideoCapturer alloc] init];
	NSLog(@"Allocating new capturer %p", capturer);
	return capturer;
}

static void releaseVideoCapturer(void *capturer){
	NSLog(@"Releasing capturer %p", capturer);
	[(VideoCapturer *)capturer stop];
	[(VideoCapturer *)capturer release];
}

static void startVideoCapture(void *capturer, uint64 opaque) {
	NSLog(@"Starting capture with capturer %p", capturer);
	[(VideoCapturer *)capturer start:(void *)opaque];
}

static void stopVideoCapture(void *capturer) {
	NSLog(@"Stopping capture with capturer %p", capturer);
	[(VideoCapturer *)capturer stop];
}

*/
import "C"
import (
	"fmt"
	"image"
	"io"
	"runtime"
	"unsafe"
)

type VideoCapturer struct {
	capturer unsafe.Pointer // Can't use C.VideoCapturer because objective-C objects aren't handled completely by Go
	frames   chan (image.Image)
}

func NewVideoCapturer() *VideoCapturer {
	c := &VideoCapturer{
		capturer: C.newVideoCapturer(),
		frames:   make(chan image.Image, 4),
	}
	runtime.SetFinalizer(c, func(c *VideoCapturer) { C.releaseVideoCapturer(c.capturer) })
	return c
}

func (c *VideoCapturer) Start() error {
	// We're doing some convolutions here to pass the VideoCapturer's pointer to the C code.
	// Go doesn't like this in general because it could lead to leaks *if* we're
	// expecting the C code to manage the associated memory. In this case, the
	// VideoCapturer is always managed by Go's garbage collector, and the corresponding C
	// object is only ever deallocated in the VideoCapturer's finalizer (which also ensures
	// that the Go object always outlives the C object, making it safe to use this value
	// as a pointer in the callback below). If this does end up causing issues, the fix is
	// to store VideoCapturer identifers in a singleton map and use those identifiers as the
	// opaque value. This will necesstitate having the calling code perform a manual release
	// on the VideoCapturer.

	// TODO: error handling
	C.startVideoCapture(c.capturer, C.uint64(uintptr(unsafe.Pointer(c))))
	return nil
}

func (c *VideoCapturer) Stop() error {
	C.stopVideoCapture(c.capturer)
	return nil
}

func (c *VideoCapturer) processFrame(img image.Image) {
	c.frames <- img
	fmt.Printf("Got image in callback: %v\n", img.Bounds())
}

func (c *VideoCapturer) FrameChannel() <-chan image.Image {
	return c.frames
}

// TODO -- remove this to a dedicated "pion adapter" class
func (c *VideoCapturer) Read() (img image.Image, release func(), err error) {
	release = func() {}
	img = <-c.frames
	if img == nil {
		err = io.EOF
	}
	return
}

//export process_yuv_frame
func process_yuv_frame(
	opaque unsafe.Pointer,
	y unsafe.Pointer,
	cb unsafe.Pointer,
	cr unsafe.Pointer,
	yStride C.int,
	cStride C.int,
	width C.int,
	height C.int,
) {
	c := (*VideoCapturer)(opaque)
	img := &image.YCbCr{
		Y:              C.GoBytes(y, yStride*height),
		Cb:             C.GoBytes(cb, cStride*height/2),
		Cr:             C.GoBytes(cr, cStride*height/2),
		YStride:        int(yStride),
		CStride:        int(cStride),
		SubsampleRatio: image.YCbCrSubsampleRatio420,
		Rect: image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{int(width), int(height)},
		},
	}
	c.processFrame(img)
}
