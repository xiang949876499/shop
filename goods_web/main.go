package main

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/satori/uuid"

	"shop/goods_web/global"
	"shop/goods_web/initialize"
	"shop/goods_web/middlewares"
	"shop/goods_web/utils/consul"
)

func main() {
	initialize.Init()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", middlewares.ValidateMobile)
	}

	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		global.Log.Panic("服务注册失败:", err.Error())
	}

	Router := initialize.Routers()
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		global.Log.Panic("启动失败:", err.Error())
	}

}
