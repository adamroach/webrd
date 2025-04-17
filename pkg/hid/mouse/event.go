package mouse

type Event struct {
	Button int  `json:"button"` // The button that was pressed or released
	Down   bool `json:"down"`   // True if the button was pressed down, false if it was released
}
