package global

import (
	"shop/order_web/config"
	"shop/order_web/proto"

	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"
)

var (
	Trans        ut.Translator
	ServerConfig config.ServerConfig
	LocalConfig  config.LocalConfig
	Log          *zap.SugaredLogger

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
	OrderSrvClient     proto.OrderClient
)
