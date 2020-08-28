package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	pb "learnrcp/grpc_01/demoproto"
	"log"
	"net"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

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
	//var opts []grpc.ServerOption
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			LoggingInterceptor,
		),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDemoProtoServer(grpcServer, &serverOp{})
	grpcServer.Serve(lis)
}