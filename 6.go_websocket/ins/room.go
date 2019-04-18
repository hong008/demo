package __go_websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
)

//房间信息，内部包含所有的客户端连接
type room interface {
	broadcast([]byte) error
	list() []client
	getClientById(int32) client
}

type myRoom struct {
	clients map[client]bool //客户端
}

func (r *myRoom) init() {
	r.clients = make(map[client]bool)
}

//广播消息
func (r *myRoom) broadcast(message []byte) error {
	for _, c := range r.list() {
		err := c.write(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *myRoom) list() []client {
	var result []client
	for c := range r.clients {
		result = append(result, c)
	}
	return result
}

func (r *myRoom) getClientById(id int32) client {
	for c := range r.clients {
		if c.getId() == id {
			return c
		}
	}
	return nil
}

func Run(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	client := &myClient{
		conn: c,
	}
}
