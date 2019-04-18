package __go_websocket

import (
	"github.com/gorilla/websocket"
)

type client interface {
	getId() int32              //获取ID
	receive()                  //接收来自消息，收到消息后的操作: 转发给其他client
	write([]byte) error        //发送消息
	sendTo(int32, []byte) bool //发送消息给指定客户端
}

type myClient struct {
	room room
	id   int32           //每个客户端对应一个ID
	conn *websocket.Conn //客户端的ws连接
}

func (c *myClient) getId() int32 {
	return c.id
}

func (c *myClient) receive() {
	//设置最大消息长度
	c.conn.SetReadLimit(512)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.room.broadcast(message)
	}
}

func (c *myClient) sendTo(id int32, message []byte) bool {
	target := c.room.getClientById(id)
	if target == nil {
		return false
	}
	if err := target.write(message); err != nil {
		return false
	}
	return true
}

func (c *myClient) write(message []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, message)
}
