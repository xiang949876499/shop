package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.32.192:9876"}))
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}

	res, err := p.SendSync(context.Background(), primitive.NewMessage("hello", []byte("this is a test")))
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("发送成功:%s\n", res.String())
	}

	if err = p.Shutdown(); err != nil {
		panic(err)
	}
}
