package grpc

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/zelas91/metric-collector/api/gen"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/types"
	"github.com/zelas91/metric-collector/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"time"
)

var (
	log = logger.New()
)

type ClientGRPC struct {
	rpc pb.MetricsClient
	IP  string
}

func HashCalcInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()

	// вызываем RPC-метод
	in, ok := req.(*pb.MetricArray)
	if ok {
		fmt.Println("OK")
	}
	b, err := proto.Marshal(in)
	if err != nil {
		log.Errorf("PZDS MARSHAL ")
	}

	fmt.Println(len(b))
	err = invoker(ctx, method, req, reply, cc, opts...)
	// выполняем действия после вызова метода\
	in, ok = req.(*pb.MetricArray)
	if ok {
		fmt.Println("OK2")
	}
	b, err = proto.Marshal(in)
	if err != nil {
		log.Errorf("PZDS MARSHAL ")
	}

	fmt.Println(len(b))
	if err != nil {
		log.Errorf("[ERROR] %s,%v", method, err)
	} else {
		log.Errorf("[INFO] %s,%v", method, time.Since(start))
	}

	return err
}

func NewClientGRPC(addr string) *ClientGRPC {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(HashCalcInterceptor),
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
func UpdateMetricsGRPC(baseURL, key string, pubKey *rsa.PublicKey, report <-chan []repository.Metric, exit <-chan time.Time) {

	client := NewClientGRPC("")
	//IP, err := getInterfaceIP("eth0")
	//if err != nil {
	//	log.Error(err)
	//} else {
	//	client.IP = IP
	//}

	for m := range report {
		headers := make(map[string]string)
		arrayMetrics, err := convertArrayMetricsToArrayMetricsGRPC(m)
		if err != nil {
			log.Errorf("error convert metrics to GRPC %v", err)
			continue
		}

		body, err := proto.Marshal(arrayMetrics)
		if err != nil {
			log.Errorf("update metrics marshal err :%v", err)
			continue
		}
		log.Info(len(body))

		//body, err = gzipCompress(body)
		//if err != nil {
		//	log.Errorf("error compress body %v", err)
		//	continue
		//}
		//
		//body, err = crypto.Encrypt(pubKey, body)
		//if err != nil {
		//	log.Errorf("encrypt err: %v", err)
		//	continue
		//}
		//hash, err := utils.GenerateHash(body, key)
		//
		//if err != nil {
		//	if !errors.Is(err, utils.ErrInvalidKey) {
		//		log.Errorf("update metrics genetate hash err:%v", err)
		//		continue
		//	}
		//	log.Errorf("Invalid hash key")
		//}
		//
		//if hash != nil {
		//	headers["HashSHA256"] = *hash
		//}
		////if client.IP != "" {
		////	headers["X-Real-IP"] = client.IP
		////}
		headers["Content-Type"] = "application/json"
		headers["Content-Encoding"] = "gzip"

		//if err = requestPost(client.client, headers, body, baseURL); err != nil {
		//	r := retryUpdateMetrics(requestPost, exit)
		//	if err = r(client.client, headers, body, baseURL); err != nil {
		//		log.Errorf("retry err: %v", err)
		//	}
		//}
		md := metadata.New(headers)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		_, _ = client.rpc.AddMetrics(ctx, arrayMetrics)
	}
}
