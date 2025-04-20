package hid

import "github.com/adamroach/webrd/pkg/hid/key"

type Keyboard interface {
	Key(event key.Event) error
}
