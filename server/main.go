package main

import (
	"log"
	"mriedel/chat/server/pkg/websocket"
	"net/http"
)

func setupRoutes() {
	room := websocket.NewRoom()
	go room.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.Serve(room, w, r)
	})
}

func main() {
	log.Printf("Starting server...")
	setupRoutes()
	log.Printf("Server started! Listening on port: 9876\n")
	log.Fatal(http.ListenAndServe(":9876", nil))

}
