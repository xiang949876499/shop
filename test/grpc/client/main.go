package main

import (
	"context"
	"fmt"
	"shop/test/grpc/proto"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	//go语言推荐的是返回一个error和一个正常的信息

	ctx := metadata.AppendToOutgoingContext(context.Background(), "test", "yes")
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Error was not a status error
			panic("解析error失败")
		}
		fmt.Println(st.Message())
		fmt.Println(st.Code())
	}
	fmt.Println(r.Message)
}
