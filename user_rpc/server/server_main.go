package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"io"
	pb "learnrcp/user_rpc/user_api_service"
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

var userMapDb = map[int32] *pb.User {
	25855: &pb.User{
		UserId:               25855,
		Mobile:               "18765906871",
		Nickname:             "litongxue",
	},
	25856: &pb.User{
		UserId:               25856,
		Mobile:               "13088889999",
		Nickname:             "letslets",
	},
	25857: &pb.User{
		UserId:               25857,
		Mobile:               "18655553333",
		Nickname:             "inter1908",
	},

}

func handleGetUserInfo(ctx context.Context, userId int32) (user *pb.User, err error) {
	var ok bool
	user, ok = userMapDb[userId]
	if !ok {
		err = fmt.Errorf("没有找到对应的用户信息:%d", userId)
		return
	}
	return
}

// 阻塞的双向六
func (t *serverOp) GetUserInfo(ctx context.Context, args *pb.Args) (resp *pb.Resp, err error) {
	resp = &pb.Resp{
		ErrorCode:            "",
		ErrorMsg:             "",
		Data:                 &pb.Data{
			User: &pb.User{},
		},
	}
	userId := args.UserId
	var user *pb.User
	user, err = handleGetUserInfo(ctx, userId)
	if err != nil {
		resp.ErrorCode = "-1"
		resp.ErrorMsg = err.Error()
		return
	}
	fmt.Println(*user)
	// 设置值
	resp.Data.User = user
	return
}

func (t *serverOp) GetUserInfoList(ctx context.Context, argsMany *pb.ArgsMany) (respMany *pb.RespMany, err error) {
	var user *pb.User
	respMany = &pb.RespMany{
		ErrorCode:            "",
		ErrorMsg:             "",
		DataMany:             &pb.DataMany{
			UserList:                 []*pb.User{},
		},
	}

	for _, userId := range argsMany.UserIdList {
		user, err = handleGetUserInfo(ctx, userId)
		if err != nil {
			return
		}
		respMany.DataMany.UserList = append(respMany.DataMany.UserList, user)
	}

	return
}

func (t *serverOp) GetuserInfoListStream01(stream pb.UserApiService_GetuserInfoListStream01Server) (err error) {
	var args *pb.Args
	respMany := pb.RespMany{
		ErrorCode:            "",
		ErrorMsg:             "",
		DataMany:             &pb.DataMany{
			UserList:                 []*pb.User{},
		},
	}
	for {
		args, err = stream.Recv()
		if err == io.EOF {
			// 处理完数据，一把发给客户端
			return stream.SendAndClose(&respMany)
		}
		userId := args.UserId
		var user *pb.User
		user, err = handleGetUserInfo(stream.Context(), userId)
		if err != nil {
			respMany.ErrorCode = "-1"
			respMany.ErrorMsg = err.Error()
			return
		}

		respMany.DataMany.UserList = append(respMany.DataMany.UserList, user)
	}
	return
}

// 服务端处理比较慢， 我们分批次将结果发送给客户端
func (t *serverOp) GetuserInfoListStream02(argsMany *pb.ArgsMany, stream pb.UserApiService_GetuserInfoListStream02Server) (err error) {
	var user *pb.User
	for _, userId := range argsMany.UserIdList {
		user, err = handleGetUserInfo(stream.Context(), userId)
		if err != nil {
			return
		}
		resp := &pb.Resp{
			ErrorCode:            "",
			ErrorMsg:             "",
			Data:                 &pb.Data{
				User:                user ,
			},
		}
		// 循环发送，客户端不用等着
		err = stream.Send(resp)
		if err != nil {
			return
		}
	}
	return

}

func (t *serverOp) GetUserInfoListStream03(stream pb.UserApiService_GetUserInfoListStream03Server) (err error) {

	var args *pb.Args
	resp := pb.Resp{
		ErrorCode:            "",
		ErrorMsg:             "",
		Data:             &pb.Data{
			User:                 &pb.User{},
		},
	}
	for {
		args, err = stream.Recv()
		if err == io.EOF {
			// 处理完数据，一把发给客户端
			return nil
		}
		userId := args.UserId
		var user *pb.User
		user, err = handleGetUserInfo(stream.Context(), userId)
		if err != nil {
			resp.ErrorCode = "-1"
			resp.ErrorMsg = err.Error()
			return
		}
		resp.Data.User = user
		stream.Send(&resp)
	}
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
	pb.RegisterUserApiServiceServer(grpcServer, &serverOp{})
	grpcServer.Serve(lis)
}
