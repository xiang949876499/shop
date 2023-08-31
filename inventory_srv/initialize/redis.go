package initialize

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"

	"shop/inventory_srv/global"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.Redis.Host, global.ServerConfig.Redis.Port),
		Password: global.ServerConfig.Redis.Password,
		DB:       global.ServerConfig.Redis.DB,
	})
	global.RedisClient = rdb

	pool := goredis.NewPool(rdb) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)
	global.Rs = rs
}
