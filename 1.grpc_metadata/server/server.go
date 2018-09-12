package server

import (
	"context"
	"demo/1.grpc_metadata/proto"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/12 22:06 
    @File: server.go    
*/

var (
	EmptyReqData = errors.New("please input correct request")
)

type DemoServer struct {
}

func NewDemoServer() *DemoServer {
	return &DemoServer{}
}

func (s *DemoServer) Single(ctx context.Context, req *proto.RequestInfo) (*proto.ResponseInfo, error) {
	fmt.Printf("do single...pid = [%v] content = [%v]\n", req.GetPid(), req.GetContent())

	//服务端接收来自上下文关联的md
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Printf("md = %+v  header = %v", md, md["client-md-key"])
	}

	//服务端发送md
	header := metadata.Pairs("server-header-key", "server-header-val")
	grpc.SendHeader(ctx, header)
	trailer := metadata.Pairs("server-trailer-key", "server-trailer-val")
	grpc.SetTrailer(ctx, trailer)

	rsp := &proto.ResponseInfo{
		Code:   1,
		Result: "hello",
	}
	return rsp, nil
}

func (s *DemoServer) Stream(modelServer proto.Model_StreamServer) error {
	server := NewServer(modelServer)
	server.Run()
	return nil
}

type GrpcServer struct {
	stream   proto.Model_StreamServer
	sendChan chan *proto.ResponseInfo
	recvChan chan *proto.RequestInfo
}

func NewServer(s proto.Model_StreamServer) *GrpcServer {
	return &GrpcServer{
		stream:   s,
		sendChan: make(chan *proto.ResponseInfo),
		recvChan: make(chan *proto.RequestInfo),
	}
}

func (s *GrpcServer) Run() {
	fmt.Println("GrpcServer Run...")
	if cap(s.recvChan) <= 0 {
		s.recvChan = make(chan *proto.RequestInfo, 100)
	}
	//这里监听流，recv到数据便输出到recvChan，然后执行handler
	go func() {
		s.recv()
	}()
	for {
		select {
		case req, ok := <-s.recvChan:
			if !ok {
				fmt.Println("ShutDown")
				return
			}
			err := s.handler(req)
			if err != nil {
				panic(fmt.Sprintf("handler panic...err = [%v]", err))
			}
		}
	}
}

func (s *GrpcServer) handler(req *proto.RequestInfo) error {
	if req == nil {
		return EmptyReqData
	}
	md, ok := metadata.FromContext(s.stream.Context())
	if ok {
		fmt.Printf("md = %+v\n", md)
	}
	fmt.Printf("do stream...pid = [%v] content = [%v]\n", req.GetPid(), req.GetContent())
	return nil
}

func (s *GrpcServer) recv() {
	for {
		data, err := s.stream.Recv()
		if err != nil {
			panic(fmt.Sprintf("server recv fail...[%v]", err))
		}

		//流逝服务端接收来自客户端的md
		md, ok := metadata.FromIncomingContext(s.stream.Context())
		if ok {
			fmt.Printf("stream server md = %+v\n", md)
		}

		response := &proto.ResponseInfo{
			Code:   10,
			Result: "succ",
		}
		s.Write(response)

		s.recvChan <- data
	}
}

func (s *GrpcServer) Write(response *proto.ResponseInfo) error {
	fmt.Println("in stream write...")
	header := metadata.Pairs("stream-server-header-key", "stream-server-header-val")
	trailer := metadata.Pairs("stream-server-trailer-key", "stream-server-trailer-val")
	s.stream.SetHeader(header)
	s.stream.SetTrailer(trailer)
	return s.stream.Send(response)
}
