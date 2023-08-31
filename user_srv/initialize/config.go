package initialize

import (
	"encoding/json"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"

	"shop/user_srv/global"
)

func InitConfig() {
	configName := "./config.yaml"
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(&global.LocalConfig); err != nil {
		panic(err)
	}
	global.Log.Infof("配置信息: %v", global.LocalConfig)

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.LocalConfig.Nacos.Host,
			Port:   global.LocalConfig.Nacos.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         global.LocalConfig.Nacos.Namespace,
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
		DataId: global.LocalConfig.Nacos.DataId,
		Group:  global.LocalConfig.Nacos.Group})

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		panic(err)
	}
	global.Log.Infof("ServerConfig :%v", global.ServerConfig)
}
