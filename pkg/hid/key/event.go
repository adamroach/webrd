package key

type Location int

const (
	LocationStandard Location = iota // Standard location
	LocationLeft                     // Left location
	LocationRight                    // Right location
	LocationNumpad                   // Numpad location
)

// See https://developer.mozilla.org/en-US/docs/Web/API/UI_Events/Keyboard_event_key_values
// for the list of key values.

type Event struct {
	Key      string   `json:"key"`      // The logical meaning of the key
	Code     Code     `json:"code"`     // The physical key code
	Location Location `json:"location"` // The location of the key on the keyboard
	KeyDown  bool     `json:"keyDown"`  // True if the key was pressed down, false if it was released
}
