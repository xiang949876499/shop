package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
	goredislib "github.com/redis/go-redis/v9"
)

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.32.192:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	gNum := 2
	var wg sync.WaitGroup
	wg.Add(gNum)
	//mutexname := "111"

	//for i := 0; i < gNum; i++ {
	//	go func() {
	//		defer wg.Done()
	//		mutex := rs.NewMutex(mutexname)
	//		fmt.Println("开始获取锁 = ", i)
	//		if err := mutex.Lock(); err != nil {
	//			panic(err)
	//		}
	//
	//		fmt.Println("获取锁成功 = ", i)
	//		// Do your work that requires the lock.
	//		time.Sleep(time.Second * 5)
	//		fmt.Println("开始释放锁 =", i)
	//		// Release the lock so other processes or threads can obtain a lock.
	//		if ok, err := mutex.Unlock(); !ok || err != nil {
	//			panic("unlock failed")
	//		}
	//	}()
	//}

	mutexname := "421"

	for i := 0; i < gNum; i++ {
		go func() {
			defer wg.Done()
			mutex := rs.NewMutex(mutexname, redsync.WithExpiry(time.Second*10))
			//zookeeper的分布式锁 -

			fmt.Println("开始获取锁")
			if err := mutex.Lock(); err != nil {
				panic(err)
			}

			fmt.Println("获取锁成功")

			time.Sleep(time.Second * 5)

			fmt.Println("开始释放锁")
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			fmt.Println("释放锁成功")
		}()
	}
	wg.Wait()
}
