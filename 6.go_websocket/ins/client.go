package ins

import (
	"github.com/gorilla/websocket"
)

type Client interface {
	GetId() int32              //获取ID
	Receive()                  //接收来自消息，收到消息后的操作: 转发给其他client
	Write([]byte) error        //发送消息
	SendTo(int32, []byte) bool //发送消息给指定客户端
}

type myClient struct {
	room Room
	id   int32           //每个客户端对应一个ID
	conn *websocket.Conn //客户端的ws连接
}

func NewClient(roon Room, conn *websocket.Conn) Client {
	return &myClient{
		room: roon,
		conn: conn,
	}
}

func (c *myClient) GetId() int32 {
	return c.id
}

func (c *myClient) Receive() {
	//设置最大消息长度
	c.conn.SetReadLimit(512)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.room.UnRegister(c)
			c.conn.Close()
			break
		}
		c.room.Broadcast(message)
	}
}

func (c *myClient) SendTo(id int32, message []byte) bool {
	target := c.room.GetClientById(id)
	if target == nil {
		return false
	}
	if err := target.Write(message); err != nil {
		return false
	}
	return true
}

func (c *myClient) Write(message []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, message)
}
