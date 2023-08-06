package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/storages"
	"net/http"
)

func Run(endpointServer string) error {
	gin.SetMode(gin.ReleaseMode)
	handler := controller.NewHandler(storages.NewMemStorage())
	return http.ListenAndServe(endpointServer, handler.InitRoutes())
}
