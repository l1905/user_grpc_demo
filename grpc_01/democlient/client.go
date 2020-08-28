package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	pb "learnrcp/grpc_01/demoproto"
	"log"
)


func main() {
	flag.Parse()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial("localhost:1024", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewDemoProtoClient(conn)

	resp, err := client.Divide(context.Background(), &pb.Args{A: 46, B: 23})

	fmt.Println("hello world")
	fmt.Println(resp.R)
	fmt.Println(err)

}