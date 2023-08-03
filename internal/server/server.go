package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/handlers"
	"github.com/zelas91/metric-collector/internal/server/storages"
	"net/http"
)

func Run(port string) error {
	gin.SetMode(gin.ReleaseMode)
	handler := handlers.NewHandler(storages.NewMemStorage())
	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler.InitRoutes())
}
