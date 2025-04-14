//go:build cgo && darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework AVFoundation -framework Foundation -framework CoreMedia -framework CoreVideo

#import "capturer.h"

static void* newCapturer(){
    return [[Capturer alloc] init];
}

static void releaseCapturer(void *capturer){
	[(Capturer *)capturer stop];
	[(Capturer *)capturer release];
}

static void startCapture(void *capturer, uint64 opaque) {
	[(Capturer *)capturer start:(void *)opaque];
}

static void stopCapture(void *capturer) {
	[(Capturer *)capturer stop];
}

*/
import "C"
import (
	"fmt"
	"image"
	"io"
	"runtime"
	"time"
	"unsafe"
)

type Capturer struct {
	capturer unsafe.Pointer // Can't use C.Capturer because objective-C objects aren't handled completely by Go
	frames   chan (image.Image)
}

func NewCapturer() *Capturer {
	c := &Capturer{
		capturer: C.newCapturer(),
		frames:   make(chan image.Image, 4),
	}
	runtime.SetFinalizer(c, func(c *Capturer) { C.releaseCapturer(c.capturer) })
	return c
}

func (c *Capturer) Start() {
	// We're doing some convolutions here to pass the Capturer's pointer to the C code.
	// Go doesn't like this in general because it could lead to leaks *if* we're
	// expecting the C code to manage the associated memory. In this case, the
	// Capturer is always managed by Go's garbage collector, and the corresponding C
	// object is only ever deallocated in the Capturer's finalizer (which also ensures
	// that the Go object always outlives the C object, making it safe to use this value
	// as a pointer in the callback below). If this does end up causing issues, the fix is
	// to store Capturer identifers in a singleton map and use those identifiers as the
	// opaque value. This will necesstitate having the calling code perform a manual release
	// on the Capturer.
	C.startCapture(c.capturer, C.uint64(uintptr(unsafe.Pointer(c))))
	time.Sleep(1500 * time.Millisecond)
	c.Stop()
}

func (c *Capturer) Stop() {
	C.stopCapture(c.capturer)
}

func (c *Capturer) processFrame(img image.Image) {
	c.frames <- img
	fmt.Printf("Got image in callback: %v\n", img.Bounds())
}

func (c *Capturer) Read() (img image.Image, release func(), err error) {
	release = func() {}
	img = <-c.frames
	if img == nil {
		err = io.EOF
	}
	return
}

//export process_frame
func process_frame(
	opaque unsafe.Pointer,
	y unsafe.Pointer,
	cb unsafe.Pointer,
	cr unsafe.Pointer,
	yStride C.int,
	cStride C.int,
	width C.int,
	height C.int,
) {
	c := (*Capturer)(opaque)
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
