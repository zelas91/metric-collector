package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/server/storages"
	"log"
	"net/http"
	"time"
)

var serv *http.Server

func Run(endpointServer string) {
	gin.SetMode(gin.ReleaseMode)
	metric := controller.NewMetricHandler(service.NewMetricsService(storages.NewMemStorage()))
	serv = &http.Server{
		Addr:    endpointServer,
		Handler: metric.InitRoutes(), // Ваш обработчик запросов
	}
	go func() {
		if err := serv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe %v", err)
		}
	}()
}
func Shutdown(ctx context.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctxTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("shutdown server %v", err)
	}
}
