package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/satori/uuid"

	"shop/user_web/global"
	"shop/user_web/initialize"
	"shop/user_web/middlewares"
	"shop/user_web/utils"
	"shop/user_web/utils/consul"
)

func main() {
	initialize.Init()

	port, _ := utils.GetFreePort()
	port = 7001 //测试
	Router := initialize.Routers()
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", middlewares.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host, port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		global.Log.Panic("服务注册失败:", err.Error())
	}

	global.Log.Debugf("启动服务器, 端口： %d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		global.Log.Panic("启动失败:", err.Error())
	}
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
