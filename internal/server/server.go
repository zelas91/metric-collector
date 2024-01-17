// Package server start and shutdown web server.

package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	pb "github.com/zelas91/metric-collector/internal/api"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/utils"
	"github.com/zelas91/metric-collector/internal/utils/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
	"time"
)

var (
	serv *Server
	log  = logger.New()
)

type Server struct {
	http     *http.Server
	repo     repository.StorageRepository
	grpcServ *grpc.Server
}

// Run start web server.
func Run(ctx context.Context, cfg *config.Config) {
	gin.SetMode(gin.ReleaseMode)

	var repo repository.StorageRepository

	if cfg.Database != nil && *cfg.Database != "" {
		repo = repository.NewDBStorage(ctx, *cfg.Database)
	}
	if repo == nil {
		repo = repository.NewMemStorage()
	}
	//if repo == nil {
	//	repo = repository.NewFileStorage(ctx, cfg)
	//}
	service := service.NewMemService(ctx, repo, cfg)
	metric := controller.NewMetricHandler(service)

	serv = &Server{
		http: &http.Server{
			Addr:    *cfg.Addr,
			Handler: metric.InitRoutes(cfg.Key, crypto.LoadPrivateKey(cfg.CryptoCertPath), cfg.TrustedSubnet), // Ваш обработчик запросов
		},
		repo: repo,
	}
	go func() {
		if err := serv.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe %v", err)
		}
	}()
	servGRPC := controller.NewServerGRPC(service)
	s, err := startGRPC(servGRPC, cfg)
	if err != nil {
		log.Errorf("start grpc server err %v", err)
	}
	serv.grpcServ = s

	log.Info("start server")
}

// Shutdown web server.
func Shutdown(ctx context.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	r, ok := serv.repo.(repository.Shutdown)
	if ok {
		if err := r.Shutdown(); err != nil {
			log.Errorf("repository shutdown err %v", err)
		}
	}

	if err := serv.http.Shutdown(ctxTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("shutdown server %v", err)
	}
	serv.grpcServ.Stop()
	log.Info("server stop")
}

func startGRPC(serv *controller.ServerGRPC, cfg *config.Config) (*grpc.Server, error) {
	listen, err := net.Listen("tcp", cfg.AddrGRPC)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile(cfg.CertPath, cfg.CryptoCertPath)
	if err != nil {
		log.Errorf("Failed to create TLS credentials %v", err)
	}
	var network *net.IPNet
	if cfg.TrustedSubnet != "" {
		if n, ok := utils.GetSubnet(cfg.TrustedSubnet); ok {
			network = n
		}
	}

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(controller.TrustedSubnet(network)))
	// регистрируем сервис
	pb.RegisterMetricsServer(s, serv)

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	return s, err

}
