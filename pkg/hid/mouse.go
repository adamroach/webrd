package hid

import "github.com/adamroach/webrd/pkg/hid/mouse"

type Mouse interface {
	Move(x, y int) error
	Button(event mouse.Event) error
}
