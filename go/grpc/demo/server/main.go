package main

import (
	"context"
	"fmt"
	"grpc/demo/pb"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type T2Server struct {
	pb.UnimplementedT2ServiceServer
}

func (*T2Server) GetMsg(ctx context.Context, req *pb.T2Request) (*pb.T2Response, error) {

	return &pb.T2Response{
		Name: req.Name,
		Detail: &pb.RespDetail{
			Age:     "11",
			Address: "beijing",
		},
	}, nil
}

func main() {

	creds, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.key")
	if err != nil {
		log.Fatalf("load ssl err: %v", err)
	}

	lis, _ := net.Listen("tcp", ":5003")
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterT2ServiceServer(s, &T2Server{})
	fmt.Println("server run success")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("start server err: %v", err)
	}
}
