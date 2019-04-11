package client

import (
	"testing"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/13 00:02 
    @File: client_test.go    
*/

var (
	url = ":13140"
)

func TestSingle(t *testing.T) {
	InitConn(url)
	Single(2, "hello world")
}

func TestGrpcClient_Write(t *testing.T) {
	InitConn(url)

	client := NewClient(NewStreamClient())
	client.Run()

	/*ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			client := NewClient(NewStreamClient())
			go func() {
				client.Run()
			}()

			sendData := &proto.RequestInfo{
				Pid:     1,
				Content: "你好",
			}
			client.Write(sendData)
			client.stream.CloseSend()
		}
	}*/
}
