package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"mriedel/chat/server/pkg/socket/event"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	return conn, err
}

func Serve(room *Room, w http.ResponseWriter, r *http.Request) {
	log.Printf("Websocket endpoint called")
	conn, err := Upgrade(w, r)
	if err != nil {
		log.Printf("Error in upgrading connection: %v\n", err)
	}
	client := &Client{
		Conn: conn,
		Room: room,
		send: make(chan event.Event),
	}
	room.Register <- client
	go client.readPump()
	go client.writePump()
}
