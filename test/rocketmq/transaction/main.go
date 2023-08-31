package main

import (
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type OrderListener struct {
}

func (o *OrderListener) ExecuteLocalTransaction(addr *primitive.Message) primitive.LocalTransactionState {
	time.Sleep(time.Second * 3)
	return primitive.UnknowState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Printf("rocketmq的消息回查")
	time.Sleep(time.Second * 15)
	return primitive.RollbackMessageState
}

func main() {

}
