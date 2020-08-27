package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	pb "learnrcp/user_rpc/user_api_service"
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
	// 建立连接
	client := pb.NewUserApiServiceClient(conn)

	// 1)发起单次请求
	resp, err := client.GetUserInfo(context.TODO(), &pb.Args{
		UserId:               25855,
	})

	fmt.Printf("getUserInfo: %+v \n", *resp.Data.User)
	if err != nil {
		fmt.Printf("getUserInfo:err %s \n", err.Error())
	}

	// 2)
	respMany, err := client.GetUserInfoList(context.TODO(), &pb.ArgsMany{
		UserIdList:           []int32{25855, 25856},
	})
	for _, user := range respMany.DataMany.UserList {
		fmt.Printf("getUserInfoList: %+v \n", *user)
		if err != nil {
			fmt.Printf("getUserInfoList:err %s \n", err.Error())
		}
	}

	// 3)
	var stream01 pb.UserApiService_GetuserInfoListStream01Client
	stream01, err = client.GetuserInfoListStream01(context.TODO())
	userIdList := []int32{25855, 25856}
	for _, userId := range userIdList {
		err = stream01.Send(&pb.Args{UserId:userId})
		if err != nil {
			fmt.Printf("GetuserInfoListStream01:err %s \n", err.Error())
		}
	}
	respMany, err = stream01.CloseAndRecv()
	if err != nil {
		fmt.Printf("GetuserInfoListStream01:err %s \n", err.Error())
	}
	for _, user := range respMany.DataMany.UserList {
		fmt.Printf("getUserInfoList01: %+v \n", *user)
		if err != nil {
			fmt.Printf("getUserInfoList01:err %s \n", err.Error())
		}
	}

	// 4)
	var stream02 pb.UserApiService_GetuserInfoListStream02Client
	stream02, err = client.GetuserInfoListStream02(context.TODO(), &pb.ArgsMany{
		UserIdList:           []int32{25855, 25856},
	})
	for {
		resp, err = stream02.Recv()
		if err == io.EOF {
			// 处理完数据，一把发给客户端
			break
		}
		fmt.Printf("GetuserInfoListStream02: %+v \n", *resp.Data.User)
	}
	stream02.Recv()


	// 5)
	var stream03 pb.UserApiService_GetUserInfoListStream03Client
	stream03, err = client.GetUserInfoListStream03(context.TODO())

	for _, userId := range []int32{25855, 25856, 25857} {
		stream03.Send(&pb.Args{UserId:userId})
	}
	for {
		// 异步接收数据
		resp, err = stream03.Recv()
		if err == io.EOF {
			// 处理完数据，一把发给客户端
			break
		}
		fmt.Printf("GetuserInfoListStream03: %+v \n", *resp.Data.User)
	}
	return


}


