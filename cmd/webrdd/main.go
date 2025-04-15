package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/adamroach/webrd/pkg/capture"
)

func main() {
	c := capture.NewVideoCapturer()
	c.Start()
	go func() {
		i := 0
		lastFrame := time.Now()
		for frame := range c.FrameChannel() {
			i++
			fmt.Printf("Read frame %d: %v (%v)\n", i, frame.Bounds(), time.Since(lastFrame))
			lastFrame = time.Now()
		}
	}()
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("Stopping capture")
	c.Stop()
	c = nil
	fmt.Println("Stopped capture")
	runtime.GC()
	time.Sleep(500 * time.Millisecond)
}
