package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	conn *websocket.Conn
	send chan []byte
	recv chan []byte
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO -- make this more secure
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err) // TODO -- send error to client
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		w.Write(fmt.Appendf(nil, "Could not upgrade connection: %v", err))
		return
	}
	client := &WebSocket{
		conn: conn,
		send: make(chan []byte, 100),
		recv: make(chan []byte, 100),
	}
	go client.readMessages()
	go client.writeMessages()
}

func (ws *WebSocket) Send(message any) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	ws.send <- jsonMessage
	return nil
}

func (ws *WebSocket) Receive() (any, error) {
	msg, ok := <-ws.recv
	if !ok {
		return nil, io.EOF
	}
	return MakeMessage(msg)
}

func (ws *WebSocket) Close() error {
	close(ws.send)
	err := ws.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (ws *WebSocket) readMessages() {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		ws.recv <- msg
	}
	close(ws.recv)
	ws.conn.Close()
}

func (ws *WebSocket) writeMessages() {
	for msg := range ws.send {
		err := ws.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
	ws.conn.Close()
}
