package ins

import "sync"

//房间信息，内部包含所有的客户端连接
type Room interface {
	Broadcast([]byte) error
	List() []Client
	GetClientById(int32) Client
	Register(client Client)
	UnRegister(client Client)
}

type myRoom struct {
	sync.RWMutex
	clients map[Client]bool //客户端
}

func NewRoom() Room {
	return &myRoom{
		clients: map[Client]bool{},
	}
}

func (r *myRoom) init() {
	r.clients = make(map[Client]bool)
}

//广播消息
func (r *myRoom) Broadcast(message []byte) error {
	for _, c := range r.List() {
		err := c.Write(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *myRoom) List() []Client {
	r.RLock()
	var result []Client
	for c := range r.clients {
		result = append(result, c)
	}
	r.RUnlock()
	return result
}

func (r *myRoom) GetClientById(id int32) Client {
	r.RLock()
	for c := range r.clients {
		if c.GetId() == id {
			r.RUnlock()
			return c
		}
	}
	r.RUnlock()
	return nil
}

func (r *myRoom) Register(client Client) {
	r.Lock()
	r.clients[client] = true
	r.Unlock()
	return
}

func (r *myRoom) UnRegister(client Client) {
	r.Lock()
	delete(r.clients, client)
	r.Unlock()
}
