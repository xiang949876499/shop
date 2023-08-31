package global

import (
	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"shop/inventory_srv/config"
)

var (
	DB           *gorm.DB
	RedisClient  *redis.Client
	Rs           *redsync.Redsync
	EsClient     *elastic.Client
	Log          *zap.SugaredLogger
	ServerConfig config.ServerConfig
	LocalConfig  config.LocalConfig
)
