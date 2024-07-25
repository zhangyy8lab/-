package main

import (
	"context"
	"fmt"
	"grpc/demo/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":5003", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client start err: %v", err)
	}

	defer conn.Close()
	c := pb.NewT2ServiceClient(conn)
	resp, err := c.GetMsg(context.Background(), &pb.T2Request{Name: "test"})
	if err != nil {
		log.Fatalf("client call err: %v", err)
	}
	fmt.Println(resp.Name, resp.Detail.Age, resp.Detail.Address)
}
