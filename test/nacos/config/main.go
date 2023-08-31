package main

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func main() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.32.192",
			Port:   8848,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         "02458aa2-cce2-4c30-905c-fc6399e03d4d",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-web",
		Group:  "test"})

	if err != nil {
		panic(err)
	}
	fmt.Printf(content)

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web",
		Group:  "test",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Printf("配置文件发生变化")
		},
	})
	select {}
}
