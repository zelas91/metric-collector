// Package server start and shutdown web server.

package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/utils/crypto"
	"net/http"
	"time"
)

var (
	serv *Server
	log  = logger.New()
)

type Server struct {
	http *http.Server
	repo repository.StorageRepository
}

// Run start web server.
func Run(ctx context.Context, cfg *config.Config) {
	gin.SetMode(gin.ReleaseMode)

	var repo repository.StorageRepository

	if cfg.Database != nil && *cfg.Database != "" {
		repo = repository.NewDBStorage(ctx, *cfg.Database)
	}
	if repo == nil && (cfg.Restore == nil || cfg.FilePath == nil) {
		repo = repository.NewMemStorage()
	}
	if repo == nil {
		repo = repository.NewFileStorage(ctx, cfg)
	}

	metric := controller.NewMetricHandler(service.NewMemService(ctx, repo, cfg))

	serv = &Server{
		http: &http.Server{
			Addr:    *cfg.Addr,
			Handler: metric.InitRoutes(cfg.Key, crypto.LoadPrivateKey(cfg.CryptoCertPath)), // Ваш обработчик запросов
		},
		repo: repo,
	}
	go func() {
		if err := serv.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe %v", err)
		}
	}()
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
	log.Info("server stop")
}
