package websocket

import (
	"log"
	"mriedel/chat/server/pkg/websocket/event"
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
			for client, _ := range r.Clients {
				err := client.WriteJSON(event.Event{
					Type: event.JoinRoomEvent,
					Data: "User joined the Room",
				})
				if err != nil {
					log.Printf("Error in writing to client %s: %v\n", client.ID, err)
					return
				}
			}
			break
		case client := <-r.Unregister:
			delete(r.Clients, client)
			for client, _ := range r.Clients {
				err := client.WriteJSON(event.Event{
					Type: event.LeaveRoomEvent,
					Data: "User left the Room",
				})
				if err != nil {
					log.Printf("Error in writing to client %s: %v\n", client.ID, err)
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
