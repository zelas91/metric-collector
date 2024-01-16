package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/zelas91/metric-collector/api/gen"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"net"
)

//func main() {
//	cfg := NewConfig()
//	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
//	_ = cancel
//	server.Run(ctx, cfg)
//	<-ctx.Done()
//	stop(ctx)
//}
//func stop(ctx context.Context) {
//	server.Shutdown(ctx)
//	logger.Shutdown()
//	os.Exit(0)
//}

type ServerGRPC struct {
	pb.UnimplementedMetricsServer
}

func (s *ServerGRPC) AddMetrics(ctx context.Context, in *pb.MetricArray) (*empty.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(md)
	}
	b, err := proto.Marshal(in)
	if err != nil {
		log.Errorf("PZDS MARSHAL ")
	}

	fmt.Println(len(b))
	return &empty.Empty{}, nil
}

func main() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServer(s, &ServerGRPC{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
