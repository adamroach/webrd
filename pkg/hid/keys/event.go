package keys

type Event struct {
	Key     rune `json:"key"`      // The key that was pressed or released; null if a control key
	Control Code `json:"control"`  // The control key that was pressed or released; null if not a control key
	KeyDown bool `json:"key_down"` // True if the key was pressed down, false if it was released
}
