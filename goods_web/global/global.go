package global

import (
	"shop/goods_web/config"
	"shop/goods_web/proto"

	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"
)

var (
	Trans          ut.Translator
	ServerConfig   config.ServerConfig
	LocalConfig    config.LocalConfig
	GoodsSrvClient proto.GoodsClient
	Log            *zap.SugaredLogger
)
