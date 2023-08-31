package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	protobuf "shop/order_srv/proto"
	"shop/order_srv/utils/consul"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"shop/order_srv/global"
	"shop/order_srv/handler"
	"shop/order_srv/initialize"
)

func main() {
	initialize.Init()
	global.Log.Info("start")

	IP := flag.String("ip", "0.0.0.0", "ip地址")

	handler.InitRocketmq()

	server := grpc.NewServer()
	protobuf.RegisterOrderServer(server, &handler.OrderServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		*IP, global.ServerConfig.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	err = registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceID)
	if err != nil {
		global.Log.Panic("服务注册失败:", err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.Deregister(serviceID); err != nil {
		global.Log.Info("注销失败")
	}
	global.Log.Info("注销成功")
}
