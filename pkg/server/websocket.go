package server

import "github.com/gorilla/websocket"

type WebSocket struct {
	conn *websocket.Conn
	send chan []byte
	recv chan []byte
}
