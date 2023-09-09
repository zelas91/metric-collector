package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/controller/middleware"
)

func (h *MetricHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(middleware.WithLogging, middleware.GzipCompress, middleware.GzipDecompress)
	router.GET("/", h.GetMetrics)
	//router.GET("/ping", h.Ping)
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
