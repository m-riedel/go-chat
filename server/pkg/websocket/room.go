package websocket

import (
	"log"
	"mriedel/chat/server/pkg/websocket/event"
	"time"
)

type Room struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan event.Event
}

func (r *Room) Start() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
			log.Printf("client joined ")
			/*
				for client, _ := range r.Clients {
					err := client.WriteJSON(event.Event{
						Type: event.JoinRoomEvent,
						Data: event.EventMessage{
							Client:    "",
							Timestamp: "",
							Message:   "User joined the room",
						},
					})
					if err != nil {
						log.Printf("Error in writing to client %s: %v\n", client.ID, err)
						return
					}
				}*/
			break
		case client := <-r.Unregister:
			delete(r.Clients, client)
			for c, _ := range r.Clients {

				err := c.WriteJSON(event.Event{
					Type: event.LeaveRoomEvent,
					Data: event.EventMessage{
						Client:    client.ID,
						Timestamp: time.Now().String(),
						Message:   "User left the room",
					},
				})
				if err != nil {
					log.Printf("Error in writing to client %s: %v\n", c.ID, err)
					return
				}
			}
			break
		case e := <-r.Broadcast:
			log.Print(e)
			for client, _ := range r.Clients {
				err := client.WriteJSON(e)
				if err != nil {
					log.Printf("Error in writing to client %s: %v\n", client.ID, err)
					return
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
