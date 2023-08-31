package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"shop/user_web/config"
	"shop/user_web/proto"
)

var (
	Trans         ut.Translator
	LocalConfig   config.LocalConfig
	UserSrvClient proto.UserClient
	ServerConfig  config.ServerConfig
	RedisClient   *redis.Client
	Log           *zap.SugaredLogger
)
