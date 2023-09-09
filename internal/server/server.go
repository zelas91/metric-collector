package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"log"
	"net/http"
	"time"
)

var serv *Server

type Server struct {
	http *http.Server
	repo repository.StorageRepository
}

func Run(ctx context.Context, cfg *config.Config) {
	gin.SetMode(gin.ReleaseMode)

	//db := repository.NewPostgresDB(*cfg.Database)

	repo := repository.NewMemStore()
	metric := controller.NewMetricHandler(service.NewMemService(ctx, repo, cfg))

	serv = &Server{
		http: &http.Server{
			Addr:    *cfg.Addr,
			Handler: metric.InitRoutes(), // Ваш обработчик запросов
		},
		repo: repo,
	}
	go func() {
		if err := serv.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe %v", err)
		}
	}()
}
func Shutdown(ctx context.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	//if err := serv.repo.Shutdown(); err != nil {
	//	log.Printf("repository shutdown err %v", err)
	//}
	if err := serv.http.Shutdown(ctxTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("shutdown server %v", err)
	}
}
