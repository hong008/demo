package client

import (
	"context"
	"demo/1.grpc_metadata/proto"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"os/signal"
	"syscall"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/12 22:10 
    @File: client.go    
*/

var (
	conn *grpc.ClientConn

	UnavailableConn = errors.New("unavailable grpc conn")
)

//新建一个Grpc连接
func InitConn(addr string) {
	var opts []grpc.DialOption
	var err error
	opts = append(opts, grpc.WithInsecure())
	conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	go sign()
}

//系统信号量监听
func sign() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			conn.Close()
			return
		}
	}
	return
}

//new metadata
func NewMD(mapData map[string]string) *metadata.MD {
	md := metadata.New(mapData)
	return &md
}

//单向模式，用的时候才调用
func Single(pid int32, content string) (*proto.ResponseInfo, error) {
	if conn == nil {
		return nil, UnavailableConn
	}
	client := proto.NewModelClient(conn)
	if client == nil {
		return nil, UnavailableConn
	}
	in := &proto.RequestInfo{
		Pid:     pid,
		Content: content,
	}
	//单向客户端给服务端发送md
	md := metadata.Pairs("client-md-key", "client-md-val")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	//rsp, err := client.Single(ctx, in)

	//单向客户端接收来自服务端的md
	header := metadata.MD{}
	trailer := metadata.MD{}
	rsp, err := client.Single(ctx, in, grpc.Header(&header), grpc.Trailer(&trailer))

	fmt.Printf("header = %+v	 trailer = %v", header, trailer)
	return rsp, err
}

//双向模式
func NewStreamClient() proto.Model_StreamClient {
	modelClient := proto.NewModelClient(conn)
	if modelClient == nil {
		fmt.Println("cannot new model client")
		return nil
	}

	//流式客户端像服务端发送md：新建一个带md的上下文
	mapParam := map[string]string{
		"user": "hong",
		"sex":  "man",
	}
	md := NewMD(mapParam)
	ctx := metadata.NewOutgoingContext(context.Background(), *md)
	client, err := modelClient.Stream(ctx)
	if err != nil {
		fmt.Printf("new stream client fail...err = [%v]\n", err)
		return nil
	}
	return client
}

/*GRPC 流式client*/
type GrpcClient struct {
	stream   proto.Model_StreamClient
	sendChan chan *proto.RequestInfo
	recvChan chan *proto.ResponseInfo
}

func NewClient(c proto.Model_StreamClient) *GrpcClient {
	return &GrpcClient{
		stream:   c,
		sendChan: make(chan *proto.RequestInfo, 100),
		recvChan: make(chan *proto.ResponseInfo, 100),
	}
}

func (c *GrpcClient) Run() {
	if cap(c.recvChan) <= 0 {
		c.recvChan = make(chan *proto.ResponseInfo, 100)
	}
	//跟server端一样，监听recv channel
	go func() {
		c.recv()
	}()

	for {
		select {
		case rsp, ok := <-c.recvChan:
			if !ok {
				panic("client run panic")
			}

			c.handler(rsp)
		}
	}
}

func (c *GrpcClient) recv() {
	for {
		data, err := c.stream.Recv()
		if err != nil {
			//panic(fmt.Sprintf("client recv fail...err = [%v]", err))
			fmt.Println(fmt.Sprintf("client recv fail...err = [%v]", err))
		}
		//流逝客户端接收来自服务端的md
		header, err := c.stream.Header()
		if err == nil {
			fmt.Printf("in client md, header = %+v\n", header)
		}
		trailer := c.stream.Trailer()
		fmt.Printf("in cliend md, trailer = %+v\n", trailer)

		c.recvChan <- data
	}
}

func (c *GrpcClient) handler(response *proto.ResponseInfo) {
	fmt.Printf("in client handler...code = [%v] result = [%v]\n", response.GetCode(), response.GetResult())
}

func (c *GrpcClient) Write(request *proto.RequestInfo) error {
	fmt.Printf("in client write...request = %+v\n", *request)
	return c.stream.Send(request)
}
