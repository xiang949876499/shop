package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"shop/user_srv/config"
)

var (
	DB           *gorm.DB
	Log          *zap.SugaredLogger
	ServerConfig config.ServerConfig
	LocalConfig  config.LocalConfig
)
