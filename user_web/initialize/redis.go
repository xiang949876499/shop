package initialize

import (
	"fmt"

	"github.com/go-redis/redis"

	"shop/user_web/global"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       global.ServerConfig.RedisInfo.DB,
	})
	global.RedisClient = rdb
}
