package global

import (
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"goods_srv/config"
)

var (
	DB           *gorm.DB
	EsClient     *elastic.Client
	Log          *zap.SugaredLogger
	ServerConfig config.ServerConfig
	LocalConfig  config.LocalConfig
)
