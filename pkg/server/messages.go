package server

import (
	"encoding/json"

	"github.com/adamroach/webrd/pkg/hid/key"
)

type MessageType string

const (
	TypeKeyboard     MessageType = "keyboard"
	TypeMouseButton  MessageType = "mouse_button"
	TypeMouseMove    MessageType = "mouse_move"
	TypeMouseWheel   MessageType = "mouse_wheel"
	TypeOffer        MessageType = "offer"
	TypeAnswer       MessageType = "answer"
	TypeIceCandidate MessageType = "candidate"
	TypeAuth         MessageType = "auth"
	TypeAuthFailure  MessageType = "auth_failure"
)

///////////////////////////////////////////////////////////////////////////
// HID messages
// These messages are sent from the client to the server to control the remote device.

type KeyboardMessage struct {
	Type  MessageType `json:"type"`
	Event key.Event   `json:"event"`
}

type MouseButtonMessage struct {
	Type   MessageType `json:"type"`
	Button int         `json:"button"`
	X      int         `json:"x"`
	Y      int         `json:"y"`
	Down   bool        `json:"down"`
}

type MouseMoveMessage struct {
	Type MessageType `json:"type"`
	X    int         `json:"x"`
	Y    int         `json:"y"`
}

type MouseWheelMessage struct {
	Type   MessageType `json:"type"`
	DeltaX int         `json:"deltaX"`
	DeltaY int         `json:"deltaY"`
	DeltaZ int         `json:"deltaZ"`
}

///////////////////////////////////////////////////////////////////////////
// WebRTC messages
// These messages are sent used to establish a WebRTC connection.

type OfferMessage struct {
	Type MessageType `json:"type"`
	SDP  string      `json:"sdp"`
}

type AnswerMessage struct {
	Type MessageType `json:"type"`
	SDP  string      `json:"sdp"`
}

type IceCandidateMessage struct {
	Type      MessageType `json:"type"`
	Candidate Candidate   `json:"candidate"`
}

type Candidate struct {
	Candidate        string `json:"candidate"`
	SdpMLineIndex    int    `json:"sdpMLineIndex"`
	SdpMid           string `json:"sdpMid"`
	UsernameFragment string `json:"usernameFragment"`
}

///////////////////////////////////////////////////////////////////////////
// Auth messages
// These messages are sent from the client to the server to authenticate the user.

type AuthMessage struct {
	Type  MessageType `json:"type"`
	Token string      `json:"token"`
}

type AuthFailureMessage struct {
	Type  MessageType `json:"type"`
	Error string      `json:"error"`
}

// /////////////////////////////////////////////////////////////////////////
func MakeMessage(bytes []byte) (msg any, err error) {
	var msgMap map[string]any
	err = json.Unmarshal(bytes, &msgMap)
	if err != nil {
		return
	}
	msgType, _ := msgMap["type"].(string)
	switch MessageType(msgType) {
	case TypeKeyboard:
		msg = &KeyboardMessage{}
	case TypeMouseButton:
		msg = &MouseButtonMessage{}
	case TypeMouseMove:
		msg = &MouseMoveMessage{}
	case TypeMouseWheel:
		msg = &MouseWheelMessage{}
	case TypeOffer:
		msg = &OfferMessage{}
	case TypeAnswer:
		msg = &AnswerMessage{}
	case TypeIceCandidate:
		msg = &IceCandidateMessage{}
	case TypeAuth:
		msg = &AuthMessage{}
	case TypeAuthFailure:
		msg = &AuthFailureMessage{}
	default:
		msg = msgMap
		return
	}
	err = json.Unmarshal(bytes, msg)
	if err != nil {
		return
	}
	return
}
