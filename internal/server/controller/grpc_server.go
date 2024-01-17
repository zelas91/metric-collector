package controller

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/zelas91/metric-collector/internal/api"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/server/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

type ServerGRPC struct {
	pb.UnimplementedMetricsServer
	memService service.Service
}

func NewServerGRPC(memService service.Service) *ServerGRPC {
	return &ServerGRPC{memService: memService}
}

func (s *ServerGRPC) AddMetrics(ctx context.Context, in *pb.MetricArray) (*empty.Empty, error) {
	metrics, err := convertArrayMetricsGRPCToArrayMetrics(in)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.memService.AddMetrics(ctx, metrics); err != nil {
		return &empty.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}
func TrustedSubnet(subnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if subnet == nil {
			return handler(ctx, req)
		}
		mb, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "get metadata error ")
		}
		var token string
		values := mb.Get("X-Real-IP")
		if len(values) > 0 {
			token = values[0]
		}
		if token != "" {
			parseIP := net.ParseIP(token)
			if parseIP == nil || !subnet.Contains(parseIP) {
				return nil, status.Error(codes.Unauthenticated, "")
			}

		}
		return handler(ctx, req)
	}
}
func convertArrayMetricsGRPCToArrayMetrics(in *pb.MetricArray) ([]repository.Metric, error) {

	metrics := make([]repository.Metric, len(in.Metrics))
	for i, val := range in.Metrics {
		tmp := repository.Metric{ID: val.Id, MType: val.MType}
		switch val.MType {
		case types.GaugeType:
			tmp.Value = &val.Value
		case types.CounterType:
			tmp.Delta = &val.Delta
		default:
			log.Info(val, tmp)
			return nil, errors.New("type error")
		}
		metrics[i] = tmp
	}
	return metrics, nil
}
