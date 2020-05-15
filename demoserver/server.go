package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "learnrcp/demoproto"
	"log"
	"net"
)

type  serverOp struct {
}

func (t *serverOp) Multiply(ctx context.Context, args *pb.Args) (resp *pb.Resp, err error) {
	resp = &pb.Resp{}
	err = nil

	resp.R = args.A * args.B

	return
}

func (t *serverOp) Divide(ctx context.Context, args *pb.Args) (resp *pb.Resp, err error) {
	resp = &pb.Resp{}
	err = nil

	resp.R = args.A / args.B
	return
}

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 1024))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDemoProtoServer(grpcServer, &serverOp{})
	grpcServer.Serve(lis)
}