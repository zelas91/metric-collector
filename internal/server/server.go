package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/storages"
	"log"
	"net/http"
)

var serv *http.Server

func Run(endpointServer string) {
	gin.SetMode(gin.ReleaseMode)
	metric := controller.NewMetricHandler(storages.NewMemStorage())
	serv = &http.Server{
		Addr:    endpointServer,
		Handler: metric.InitRoutes(), // Ваш обработчик запросов
	}
	go func() {
		err := serv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
func Shutdown(ctx context.Context) {
	err := serv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
