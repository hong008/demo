package main

import (
	"demo/6.go_websocket/ins"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	addr = "192.168.1.192:8900"
)

func main() {
	room := ins.NewRoom()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		client := ins.NewClient(room, c)
		room.Register(client)
		go client.Receive()
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
