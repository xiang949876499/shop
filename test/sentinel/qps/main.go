package main

import (
	"fmt"
	"log"

	"github.com/alibaba/sentinel-golang/core/base"

	"github.com/alibaba/sentinel-golang/core/flow"

	sentinel "github.com/alibaba/sentinel-golang/api"
)

func main() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	for i := 0; i < 12; i++ {
		e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			fmt.Println("限流了")
		} else {
			fmt.Println(e)
		}
	}

}
