package socket

import (
	"log"
	"mriedel/chat/server/pkg/socket/event"
	"time"
)

type Room struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan event.Event
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				for c, _ := range r.Clients {
					err := c.WriteJSON(event.Event{
						Type: event.LeaveRoomEvent,
						Data: event.EventMessage{
							Client:    client.Name,
							Timestamp: time.Now().String(),
							Message:   "User left the room",
						},
					})
					if err != nil {
						log.Printf("Error in writing to client %s: %v\n", c.Name, err)
						return
					}
				}
				close(client.send)
			}
		case e := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.send <- e:
				default:
					close(client.send)
					delete(r.Clients, client)
				}
			}
		}
	}
}

func NewRoom() *Room {
	return &Room{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan event.Event),
	}
}
