package hid

import "github.com/adamroach/webrd/pkg/hid/keys"

type Keyboard interface {
	Key(event keys.Event) error
}
