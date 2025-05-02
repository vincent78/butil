package main

import (
	"context"
	pb "gitee.com/vincent78/gcutil/net/gcrpc/proto.pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial(":50000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	defer conn.Close()

	c := pb.NewHelloWorldClient(conn)

	req := pb.HelloRequest{
		Name: "test",
	}

	rep, err := c.SayHello(context.Background(), &req)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	log.Print(rep.Message)
}
