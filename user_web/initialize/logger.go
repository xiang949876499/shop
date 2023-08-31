package initialize

import (
	"shop/user_web/global"

	"go.uber.org/zap"
)

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"/opt/logs/info.log",
		"stdout",
	}
	cfg.ErrorOutputPaths = []string{
		"/opt/logs/error.log",
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	//logger, _ := zap.NewDevelopment() //开发环境
	//logger, _ := zap.NewProduction() //生产环境
	zap.ReplaceGlobals(logger)
	global.Log = zap.S()
}
