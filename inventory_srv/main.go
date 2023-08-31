package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	protobuf "shop/inventory_srv/proto"
	"shop/inventory_srv/utils/consul"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"shop/inventory_srv/global"
	"shop/inventory_srv/handler"
	"shop/inventory_srv/initialize"
)

func main() {
	initialize.Init()
	global.Log.Info("start")

	IP := flag.String("ip", "0.0.0.0", "ip地址")

	server := grpc.NewServer()
	protobuf.RegisterInventoryServer(server, &handler.InventoryServer{})
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

	adds := fmt.Sprintf("%s:%d", global.ServerConfig.Rocketmq.Host, global.ServerConfig.Rocketmq.Port)
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{adds}),
		consumer.WithGroupName(global.ServerConfig.Rocketmq.GroupName),
	)

	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	if err = registerClient.Deregister(serviceID); err != nil {
		global.Log.Info("注销失败")
	}
	global.Log.Info("注销成功")
}
