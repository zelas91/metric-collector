package controller

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/controller/middleware"
	"github.com/zelas91/metric-collector/internal/utils"
)

func (h *MetricHandler) InitRoutes(hashKey *string, key *rsa.PrivateKey, subnet string) *gin.Engine {
	router := gin.New()
	if subnet != "" {
		if network, ok := utils.GetSubnet(subnet); ok {
			router.Use(middleware.TrustedSubnet(network))
		}
	}

	router.Use(middleware.HashCheck(hashKey), middleware.WithLogging, middleware.Decrypt(key),
		middleware.GzipCompress, middleware.GzipDecompress, middleware.Timeout, middleware.CalculateHash(hashKey))
	router.GET("/", h.GetMetrics)
	router.GET("/ping", h.Ping)
	router.POST("/updates", h.AddMetrics)
	update := router.Group("/update")
	value := router.Group("/value")
	{
		update.POST("/:type/:name/:value", h.AddMetric)
		update.POST("/", h.AddMetricJSON)
		value.GET("/:type/:name", h.GetMetric)
		value.POST("/", h.GetMetricJSON)
	}
	return router
}
