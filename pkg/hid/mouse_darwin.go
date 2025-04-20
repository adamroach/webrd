//go:build cgo && darwin

package hid

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics

#import <CoreGraphics/CoreGraphics.h>
#import <Carbon/Carbon.h>

void SendMouseMoveEvent(int x, int y) {
	CGEventRef event = CGEventCreateMouseEvent(NULL, kCGEventMouseMoved, CGPointMake(x, y), 0);
	CGEventPost(kCGSessionEventTap, event);
	CFRelease(event);
}

void SendMouseButtonEvent(int button, int x, int y, bool down) {
	// button: 0 = left, 1 = right, 2 = middle
	CGEventType type;
	switch (button) {
	case 0:
		button = kCGMouseButtonLeft;
		type = down ? kCGEventLeftMouseDown : kCGEventLeftMouseUp;
		break;
	case 1:
		button = kCGMouseButtonRight;
		type = down ? kCGEventRightMouseDown : kCGEventRightMouseUp;
		break;
	case 2:
		button = kCGMouseButtonCenter;
		type = down ? kCGEventOtherMouseDown : kCGEventOtherMouseUp;
		break;
	default:
		type = down ? kCGEventOtherMouseDown : kCGEventOtherMouseUp;
	}

	CGEventRef event = CGEventCreateMouseEvent(NULL, type, CGPointMake(x, y), button);
	CGEventSetIntegerValueField(event, kCGMouseEventButtonNumber, button);
	CGEventPost(kCGSessionEventTap, event);
	CFRelease(event);
}

void SendScrollWheelEvent(int deltaX, int deltaY, int deltaZ) {
	CGEventRef event = CGEventCreateScrollWheelEvent(NULL, kCGScrollEventUnitPixel, 3, deltaY, deltaX, deltaZ);
	CGEventPost(kCGSessionEventTap, event);
	CFRelease(event);
}
*/
import "C"
import "log"

/*
	void SendTouchEvent(int x, int y, int radiusX, int radiusY, int angle, float pressure) {
		// TODO -- this requires more research
		CGEventRef event = CGEventTapCreate(
			kCGSessionEventTap,
			kCGTailAppendEventTap,
			kCGEventTapOptionDefault,
			kCGEventMaskForAllEvents,
			(void *)SendTouchEvent, // This is almost certainly wrong
			NULL
		);
		CGEventSetIntegerValueField(event, kCGMouseEventButtonNumber, 0);
		CGEventSetIntegerValueField(event, kCGMouseEventPressure, pressure);
		CGEventPost(kCGSessionEventTap, event);
		CFRelease(event);
	}
*/

type darwinMouse struct{}

func NewMouse() (Mouse, error) {
	return &darwinMouse{}, nil
}

func (m *darwinMouse) Move(x, y int) error {
	C.SendMouseMoveEvent(C.int(x), C.int(y))
	return nil
}
func (m *darwinMouse) Button(button int, x int, y int, down bool) error {
	log.Printf("Button: %d, x: %d, y: %d, down: %v", button, x, y, down)
	C.SendMouseButtonEvent(C.int(button), C.int(x), C.int(y), C.bool(down))
	return nil
}
func (m *darwinMouse) Wheel(deltaX, deltaY, deltaZ int) error {
	C.SendScrollWheelEvent(-C.int(deltaX), -C.int(deltaY), -C.int(deltaZ))
	return nil
}
func (m *darwinMouse) Touch(touches []Touch, event TouchEvent) error {
	// TODO -- need to figure out how Core Graphics handles touchpad events
	return nil
}
