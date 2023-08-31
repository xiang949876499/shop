package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

func main() {
	sig := make(chan os.Signal)

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("hello"),
		//consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.32.192:9876"})),
		consumer.WithNameServer([]string{"192.168.32.192:9876"}),
	)

	err := c.Subscribe("hello", consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range ext {
				fmt.Printf("获取到的值: %v \n", ext[i])
			}
			return consumer.ConsumeSuccess, nil
		})

	_ = c.Start()

	<-sig
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}
