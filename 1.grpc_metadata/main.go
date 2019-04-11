package main

import (
	"demo/1.grpc_metadata/proto"
	"demo/1.grpc_metadata/server"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/12 23:32 
    @File: blog.go
*/

var (
	addr = ":13140"
)

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("listen fail...err = [%v]", err))
	}

	s := grpc.NewServer(grpc.MaxConcurrentStreams(1000))
	proto.RegisterModelServer(s, server.NewDemoServer())
	err = s.Serve(listener)
	if err != nil {
		panic(fmt.Sprintf("begin grpc server fail...err = [%v]", err))
	}
}
