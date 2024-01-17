package grpc

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	pb "github.com/zelas91/metric-collector/api/gen"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/types"
	"github.com/zelas91/metric-collector/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"os"
)

var (
	log = logger.New()
)

type ClientGRPC struct {
	rpc pb.MetricsClient
	IP  string
}

func getCredential(certPath string) (credentials.TransportCredentials, error) {
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); ok {
		creds := credentials.NewClientTLSFromCert(caCertPool, "")
		return creds, nil
	}
	return nil, errors.New("get cert error")
}
func NewClientGRPC(baseURL, certPath string) *ClientGRPC {
	creds, err := getCredential(certPath)
	if err != nil {
		log.Errorf("get cert %v", err)
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.Dial(baseURL, grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		log.Fatal(err)
	}
	IP, err := utils.GetInterfaceIP("eth0")
	if err != nil {
		log.Error(err)
		return &ClientGRPC{IP: "", rpc: pb.NewMetricsClient(conn)}
	}

	return &ClientGRPC{IP: IP, rpc: pb.NewMetricsClient(conn)}
}

func convertArrayMetricsToArrayMetricsGRPC(metrics []repository.Metric) (*pb.MetricArray, error) {
	metricsGRPC := make([]*pb.Metric, len(metrics))
	for i, val := range metrics {
		tmp := pb.Metric{Id: val.ID, MType: val.MType}
		switch val.MType {
		case types.GaugeType:
			tmp.Value = *val.Value
		case types.CounterType:
			tmp.Delta = *val.Delta
		default:
			return nil, errors.New("type metric error")
		}
		metricsGRPC[i] = &tmp
	}
	return &pb.MetricArray{
		Metrics: metricsGRPC,
	}, nil
}
func UpdateMetricsGRPC(baseURL, certPath string, report <-chan []repository.Metric) {
	client := NewClientGRPC(baseURL, certPath)
	IP, err := utils.GetInterfaceIP("eth0")
	if err != nil {
		log.Error(err)
	} else {
		client.IP = IP
	}

	for m := range report {

		arrayMetrics, err := convertArrayMetricsToArrayMetricsGRPC(m)
		if err != nil {
			log.Errorf("error convert metrics to GRPC %v", err)
			continue
		}
		headers := make(map[string]string)
		if client.IP != "" {
			headers["X-Real-IP"] = client.IP
		}
		md := metadata.New(headers)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		if _, err = client.rpc.AddMetrics(ctx, arrayMetrics); err != nil {
			log.Errorf("add metrics err %v", err)
		}

	}
}
