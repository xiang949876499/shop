package main

import (
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.32.192:6831",
		},
		ServiceName: "jaeger_test",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	pSpan := tracer.StartSpan("main")

	span := tracer.StartSpan("funcA", opentracing.ChildOf(pSpan.Context()))
	time.Sleep(time.Millisecond * 500)
	span.Finish()

	time.Sleep(time.Millisecond * 100)
	span2 := tracer.StartSpan("funcB", opentracing.ChildOf(pSpan.Context()))
	time.Sleep(time.Millisecond * 1000)
	span2.Finish()
	pSpan.Finish()
}
