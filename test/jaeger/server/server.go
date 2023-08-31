package main

import (
	"context"
	"net"
	"shop/test/jaeger/proto"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply,
	error) {
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func main() {
	//cfg := jaegercfg.Configuration{
	//	Sampler: &jaegercfg.SamplerConfig{
	//		Type:  jaeger.SamplerTypeConst,
	//		Param: 1,
	//	},
	//	Reporter: &jaegercfg.ReporterConfig{
	//		LogSpans:           true,
	//		LocalAgentHostPort: "192.168.32.192:6831",
	//	},
	//	ServiceName: "jaeger_test",
	//}
	//
	//tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	//if err != nil {
	//	panic(err)
	//}

	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
