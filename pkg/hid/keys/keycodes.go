package keys

type Code uint16

const (
	// ASCII where applicable
	CodeNull      Code = 0x00
	CodeBackspace Code = 0x08
	CodeTab       Code = 0x09
	CodeEnter     Code = 0x0D
	CodeEscape    Code = 0x1B
	CodeDelete    Code = 0x7F

	// Control keys
	CodeLeftShift Code = iota + 0x100
	CodeRightShift
	CodeLeftControl
	CodeRightControl
	CodeLeftAlt
	CodeRightAlt
	CodeLeftMeta
	CodeRightMeta
	CodeLeftSuper
	CodeRightSuper
	CodeCapsLock
	CodeNumLock
	CodeScrollLock

	// Special keys
	CodePause Code = iota + 0x200
	CodeInsert
	CodeHome
	CodeEnd
	CodePageUp
	CodePageDown
	CodeArrowUp
	CodeArrowDown
	CodeArrowLeft
	CodeArrowRight
	CodePrintScreen
	CodeClear

	// Function keys
	CodeF1 Code = iota + 0x300
	CodeF2
	CodeF3
	CodeF4
	CodeF5
	CodeF6
	CodeF7
	CodeF8
	CodeF9
	CodeF10
	CodeF11
	CodeF12
	CodeF13
	CodeF14
	CodeF15
	CodeF16
	CodeF17
	CodeF18
	CodeF19
	CodeF20
	CodeF21
	CodeF22
	CodeF23
	CodeF24

	// Numpad keys
	CodeNumpad0 Code = iota + 0x400
	CodeNumpad1
	CodeNumpad2
	CodeNumpad3
	CodeNumpad4
	CodeNumpad5
	CodeNumpad6
	CodeNumpad7
	CodeNumpad8
	CodeNumpad9
	CodeNumpadAdd
	CodeNumpadSubtract
	CodeNumpadMultiply
	CodeNumpadDivide
	CodeNumpadDecimal
	CodeNumpadEnter
)
