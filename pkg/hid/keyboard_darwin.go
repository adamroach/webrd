//go:build cgo && darwin

package hid

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics

#import <CoreGraphics/CoreGraphics.h>
#import <Carbon/Carbon.h>

void SendKeyEvent(CGKeyCode keyCode, bool keyDown) {
	CGEventRef event = CGEventCreateKeyboardEvent(NULL, keyCode, keyDown);
	CGEventPost(kCGSessionEventTap, event);
	CFRelease(event);
}

*/
import "C"
import (
	"log"

	"github.com/adamroach/webrd/pkg/hid/key"
)

type darwinKeyboard struct{}

func NewKeyboard() (Keyboard, error) {
	return &darwinKeyboard{}, nil
}

func (k *darwinKeyboard) Key(event key.Event) error {
	keyCode := ConvertKeyCodeToCGKeyCode(event.Code)
	log.Printf("Sending key event: %+v (%d)", event, keyCode)
	C.SendKeyEvent(keyCode, C.bool(event.KeyDown))
	return nil
}

var keymap = map[key.Code]C.CGKeyCode{
	key.CodeKeyA:            C.kVK_ANSI_A,
	key.CodeKeyS:            C.kVK_ANSI_S,
	key.CodeKeyD:            C.kVK_ANSI_D,
	key.CodeKeyF:            C.kVK_ANSI_F,
	key.CodeKeyH:            C.kVK_ANSI_H,
	key.CodeKeyG:            C.kVK_ANSI_G,
	key.CodeKeyZ:            C.kVK_ANSI_Z,
	key.CodeKeyX:            C.kVK_ANSI_X,
	key.CodeKeyC:            C.kVK_ANSI_C,
	key.CodeKeyV:            C.kVK_ANSI_V,
	key.CodeIntlBackslash:   C.kVK_ISO_Section,
	key.CodeKeyB:            C.kVK_ANSI_B,
	key.CodeKeyQ:            C.kVK_ANSI_Q,
	key.CodeKeyW:            C.kVK_ANSI_W,
	key.CodeKeyE:            C.kVK_ANSI_E,
	key.CodeKeyR:            C.kVK_ANSI_R,
	key.CodeKeyY:            C.kVK_ANSI_Y,
	key.CodeKeyT:            C.kVK_ANSI_T,
	key.CodeDigit1:          C.kVK_ANSI_1,
	key.CodeDigit2:          C.kVK_ANSI_2,
	key.CodeDigit3:          C.kVK_ANSI_3,
	key.CodeDigit4:          C.kVK_ANSI_4,
	key.CodeDigit6:          C.kVK_ANSI_6,
	key.CodeDigit5:          C.kVK_ANSI_5,
	key.CodeEqual:           C.kVK_ANSI_Equal,
	key.CodeDigit9:          C.kVK_ANSI_9,
	key.CodeDigit7:          C.kVK_ANSI_7,
	key.CodeMinus:           C.kVK_ANSI_Minus,
	key.CodeDigit8:          C.kVK_ANSI_8,
	key.CodeDigit0:          C.kVK_ANSI_0,
	key.CodeBracketRight:    C.kVK_ANSI_RightBracket,
	key.CodeKeyO:            C.kVK_ANSI_O,
	key.CodeKeyU:            C.kVK_ANSI_U,
	key.CodeBracketLeft:     C.kVK_ANSI_LeftBracket,
	key.CodeKeyI:            C.kVK_ANSI_I,
	key.CodeKeyP:            C.kVK_ANSI_P,
	key.CodeEnter:           C.kVK_Return,
	key.CodeKeyL:            C.kVK_ANSI_L,
	key.CodeKeyJ:            C.kVK_ANSI_J,
	key.CodeQuote:           C.kVK_ANSI_Quote,
	key.CodeKeyK:            C.kVK_ANSI_K,
	key.CodeSemicolon:       C.kVK_ANSI_Semicolon,
	key.CodeBackslash:       C.kVK_ANSI_Backslash,
	key.CodeComma:           C.kVK_ANSI_Comma,
	key.CodeSlash:           C.kVK_ANSI_Slash,
	key.CodeKeyN:            C.kVK_ANSI_N,
	key.CodeKeyM:            C.kVK_ANSI_M,
	key.CodePeriod:          C.kVK_ANSI_Period,
	key.CodeTab:             C.kVK_Tab,
	key.CodeSpace:           C.kVK_Space,
	key.CodeBackquote:       C.kVK_ANSI_Grave,
	key.CodeBackspace:       C.kVK_Delete,
	key.CodeEscape:          C.kVK_Escape,
	key.CodeMetaRight:       0x36,
	key.CodeMetaLeft:        C.kVK_Command,
	key.CodeShiftLeft:       C.kVK_Shift,
	key.CodeCapsLock:        C.kVK_CapsLock,
	key.CodeAltLeft:         C.kVK_Option,
	key.CodeControlLeft:     C.kVK_Control,
	key.CodeShiftRight:      C.kVK_RightShift,
	key.CodeAltRight:        C.kVK_RightOption,
	key.CodeControlRight:    C.kVK_RightControl,
	key.CodeFn:              C.kVK_Function,
	key.CodeF17:             C.kVK_F17,
	key.CodeNumpadDecimal:   C.kVK_ANSI_KeypadDecimal,
	key.CodeNumpadMultiply:  C.kVK_ANSI_KeypadMultiply,
	key.CodeNumpadAdd:       C.kVK_ANSI_KeypadPlus,
	key.CodeNumLock:         C.kVK_ANSI_KeypadClear,
	key.CodeVolumeUp:        C.kVK_VolumeUp,
	key.CodeVolumeDown:      C.kVK_VolumeDown,
	key.CodeVolumeMute:      C.kVK_Mute,
	key.CodeAudioVolumeUp:   C.kVK_VolumeUp,
	key.CodeAudioVolumeDown: C.kVK_VolumeDown,
	key.CodeAudioVolumeMute: C.kVK_Mute,
	key.CodeNumpadDivide:    C.kVK_ANSI_KeypadDivide,
	key.CodeNumpadEnter:     0x34,
	key.CodeNumpadSubtract:  C.kVK_ANSI_KeypadMinus,
	key.CodeF18:             C.kVK_F18,
	key.CodeF19:             C.kVK_F19,
	key.CodeNumpadEqual:     C.kVK_ANSI_KeypadEquals,
	key.CodeNumpad0:         C.kVK_ANSI_Keypad0,
	key.CodeNumpad1:         C.kVK_ANSI_Keypad1,
	key.CodeNumpad2:         C.kVK_ANSI_Keypad2,
	key.CodeNumpad3:         C.kVK_ANSI_Keypad3,
	key.CodeNumpad4:         C.kVK_ANSI_Keypad4,
	key.CodeNumpad5:         C.kVK_ANSI_Keypad5,
	key.CodeNumpad6:         C.kVK_ANSI_Keypad6,
	key.CodeNumpad7:         C.kVK_ANSI_Keypad7,
	key.CodeF20:             C.kVK_F20,
	key.CodeNumpad8:         C.kVK_ANSI_Keypad8,
	key.CodeNumpad9:         C.kVK_ANSI_Keypad9,
	key.CodeIntlYen:         C.kVK_JIS_Yen,
	key.CodeIntlRo:          C.kVK_JIS_Underscore,
	key.CodeNumpadComma:     C.kVK_JIS_KeypadComma,
	key.CodeF5:              C.kVK_F5,
	key.CodeF6:              C.kVK_F6,
	key.CodeF7:              C.kVK_F7,
	key.CodeF3:              C.kVK_F3,
	key.CodeF8:              C.kVK_F8,
	key.CodeF9:              C.kVK_F9,
	key.CodeLang2:           C.kVK_JIS_Eisu,
	key.CodeF11:             C.kVK_F11,
	key.CodeLang1:           C.kVK_JIS_Kana,
	key.CodeF13:             C.kVK_F13,
	key.CodeF16:             C.kVK_F16,
	key.CodeF14:             C.kVK_F14,
	key.CodeF10:             C.kVK_F10,
	key.CodeContextMenu:     0x6E,
	key.CodeF12:             C.kVK_F12,
	key.CodeF15:             C.kVK_F15,
	key.CodeHelp:            C.kVK_Help,
	key.CodeHome:            C.kVK_Home,
	key.CodePageUp:          C.kVK_PageUp,
	key.CodeDelete:          C.kVK_ForwardDelete,
	key.CodeF4:              C.kVK_F4,
	key.CodeEnd:             C.kVK_End,
	key.CodeF2:              C.kVK_F2,
	key.CodePageDown:        C.kVK_PageDown,
	key.CodeF1:              C.kVK_F1,
	key.CodeArrowLeft:       C.kVK_LeftArrow,
	key.CodeArrowRight:      C.kVK_RightArrow,
	key.CodeArrowDown:       C.kVK_DownArrow,
	key.CodeArrowUp:         C.kVK_UpArrow,
}

func ConvertKeyCodeToCGKeyCode(keyCode key.Code) C.CGKeyCode {
	code, ok := keymap[keyCode]
	if ok {
		return code
	}
	return C.CGKeyCode(0)
}
