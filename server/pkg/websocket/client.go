package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"mriedel/chat/server/pkg/websocket/event"
	"sync"
	"time"
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
	/*
		c.Conn.SetCloseHandler(func(code int, text string) error {
			c.Room.Unregister <- c
			return c.Conn.Close()
		})*/
	for {
		var e event.Event
		err := c.Conn.ReadJSON(&e)
		if err != nil {
			log.Printf("Error in Reading the message: %v\n", err)
			return
		}
		switch e.Type {
		//case event.LeaveRoomEvent:
		//case event.JoinRoomEvent:
		case event.SetUsernameEvent:
			c.ID = e.Data.Client
			c.Room.Broadcast <- event.Event{
				Type: event.JoinRoomEvent,
				Data: event.EventMessage{
					Client:    c.ID,
					Timestamp: time.Now().String(),
					Message:   "User joined the room!",
				},
			}
		case event.MessageEvent:
			c.Room.Broadcast <- event.Event{
				Type: event.MessageEvent,
				Data: event.EventMessage{
					Client:    c.ID,
					Timestamp: time.Now().String(),
					Message:   e.Data.Message,
				},
			}
		}
	}
}

func (c *Client) WriteJSON(e event.Event) error {
	return c.Conn.WriteJSON(e)
}
