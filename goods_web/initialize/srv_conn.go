package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" //It's important 不引入会引起 too many colons in address 报错
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"shop/goods_web/global"
	"shop/goods_web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.Consul
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodSrv.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		global.Log.Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	global.GoodsSrvClient = proto.NewGoodsClient(userConn)
}
