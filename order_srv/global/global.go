package global

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"shop/order_srv/config"
	"shop/order_srv/proto"
)

var (
	DB                 *gorm.DB
	EsClient           *elastic.Client
	Rocketmq           rocketmq.TransactionProducer
	Log                *zap.SugaredLogger
	ServerConfig       config.ServerConfig
	LocalConfig        config.LocalConfig
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
