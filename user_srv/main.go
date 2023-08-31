package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"shop/user_srv/global"
	"shop/user_srv/headler"
	"shop/user_srv/initialize"
	"shop/user_srv/proto"
	"shop/user_srv/utils"
	"shop/user_srv/utils/consul"
)

func main() {
	initialize.Init()
	global.Log.Info("start ", global.ServerConfig)

	port, _ := utils.GetFreePort()

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &headler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		global.ServerConfig.Host, port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	err = registerClient.Register(global.ServerConfig.Host, port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceID)
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
