package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"mriedel/chat/server/pkg/socket/event"
	"time"
)

type Client struct {
	Name string
	Conn *websocket.Conn
	Room *Room
	send chan event.Event
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()
	//c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var e event.Event
		err := c.Conn.ReadJSON(&e)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		switch e.Type {
		case event.SetUsernameEvent:
			c.Name = e.Data.Client
			c.Room.Broadcast <- event.Event{
				Type: event.JoinRoomEvent,
				Data: event.EventMessage{
					Client:    c.Name,
					Timestamp: time.Now().String(),
					Message:   "User joined the room!",
				},
			}
		case event.MessageEvent:
			c.Room.Broadcast <- event.Event{
				Type: event.MessageEvent,
				Data: event.EventMessage{
					Client:    c.Name,
					Timestamp: time.Now().String(),
					Message:   e.Data.Message,
				},
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case e, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteJSON(e)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) WriteJSON(e event.Event) error {
	return c.Conn.WriteJSON(e)
}
