package initialize

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func TestRedis(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "192.168.32.192", 6379),
		Password: "",
		DB:       0,
	})
	fmt.Println(rdb.Ping())

}
