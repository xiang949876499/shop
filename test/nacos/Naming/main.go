package main

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func main() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "02458aa2-cce2-4c30-905c-fc6399e03d4d", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.32.192",
			Port:   8848,
		},
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	//注册实例
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "192.168.32.192",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		//ClusterName: "users", // 默认值DEFAULT
		//GroupName:   "users", // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("注册服务 %v", success)

	services, err := namingClient.GetService(vo.GetServiceParam{
		ServiceName: "demo.go",
		//Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
		//GroupName:   "group-a",             // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("获取服务 = %v", services)

	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	instances, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		//GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		//Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("实列列表 = %v ", instances)

	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
		//GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		//Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
	})

	fmt.Printf("轮询列表 = %v ", instance)

	serviceInfos, err := namingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: "02458aa2-cce2-4c30-905c-fc6399e03d4d",
		PageNo:    1,
		PageSize:  10,
	})
	fmt.Printf("获取服务名列表 = %v ", serviceInfos)

	select {}
}
