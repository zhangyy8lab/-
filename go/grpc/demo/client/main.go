package main

import (
	"context"
	"fmt"
	"grpc/demo/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 以下是使用ssl
	creds, err := credentials.NewClientTLSFromFile("ssl/server.crt", "grpc.demo.test")
	if err != nil {
		log.Fatalf("client load ssl err: %v", err)
	}

	// conn, err := grpc.Dial(":5003", grpc.WithInsecure()) // 不使用ssl
	conn, err := grpc.Dial(":5003", grpc.WithTransportCredentials(creds)) // 使用ssl
	// conn, err := grpc.Dial(":5003", grpc.WithTransportCredentials(creds)) // 使用ssl
	if err != nil {
		log.Fatalf("client start err: %v", err)
	}

	defer conn.Close()
	fmt.Println("aa")
	c := pb.NewT2ServiceClient(conn)
	resp, err := c.GetMsg(context.Background(), &pb.T2Request{Name: "test"})
	if err != nil {
		log.Fatalf("client call err: %v", err)
	}
	fmt.Println(resp.Name, resp.Detail.Age, resp.Detail.Address)
}
