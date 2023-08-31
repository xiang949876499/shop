package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" //It's important 不引入会引起 too many colons in address 报错
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"shop/order_srv/global"
	"shop/order_srv/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.Consul
	GoodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodSrv.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		global.Log.Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}

	global.GoodsSrvClient = proto.NewGoodsClient(GoodsConn)

	InventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrv.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		global.Log.Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}

	global.InventorySrvClient = proto.NewInventoryClient(InventoryConn)
}
