package hid

type Touch struct {
	X        int
	Y        int
	RadiusX  int
	RadiusY  int
	Angle    int
	Pressure float32
}

type TouchEvent string

const (
	TouchStart  TouchEvent = "touchstart"
	TouchMove   TouchEvent = "touchmove"
	TouchEnd    TouchEvent = "touchend"
	TouchCancel TouchEvent = "touchcancel"
)

type Mouse interface {
	Move(x, y int) error                              // mousemove
	Button(button int, x int, y int, down bool) error // mousedown / mouseup
	Wheel(deltaX, deltaY, deltaZ int) error           // wheel
	Touch(touches []Touch, event TouchEvent) error    // touchstart / touchmove / touchend / touchcancel
}
