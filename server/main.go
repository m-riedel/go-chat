package main

import (
	"log"
	"mriedel/chat/server/pkg/socket"
	"net/http"
)

func setupRoutes() {
	room := socket.NewRoom()
	go room.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.Serve(room, w, r)
	})
}

func main() {
	log.Printf("Starting server...")
	setupRoutes()
	log.Printf("Server started! Listening on port: 9876\n")
	log.Fatal(http.ListenAndServe(":9876", nil))

}
