package client

import (
	"demo/1.grpc_metadata/proto"
	"testing"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/13 00:02 
    @File: client_test.go    
*/

var (
	url = ":13141"
)

func TestSingle(t *testing.T) {
	InitConn(url)
	Single(2, "hello world")
}

func TestGrpcClient_Write(t *testing.T) {
	InitConn(url)
	client := NewClient(NewStreamClient())
	go func() {
		client.Run()
	}()

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			sendData := &proto.RequestInfo{
				Pid:     1,
				Content: "你好",
			}
			client.Write(sendData)
		}
	}
}
