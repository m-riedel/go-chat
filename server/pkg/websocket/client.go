package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"mriedel/chat/server/pkg/websocket/event"
	"sync"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Room *Room
	mu   sync.Mutex
}

func (c *Client) Read() {
	defer func() {
		c.Room.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		var e event.Event
		err := c.Conn.ReadJSON(&e)
		if err != nil {
			log.Printf("Error in Reading the message: %v\n", err)
			return
		}
		c.Room.Broadcast <- e
	}
}

func (c *Client) WriteJSON(e event.Event) error {
	return c.Conn.WriteJSON(e)
}
