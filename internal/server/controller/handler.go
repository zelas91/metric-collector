package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/middleware"
)

func (h *MetricHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(middleware.SetContextHTML, middleware.WithLogging)
	router.GET("/", h.GetMetrics)
	update := router.Group("/update")
	value := router.Group("/value")
	{
		update.Use(middleware.SetContextPlain)
		update.POST("/:type/:name/:value", h.AddMetric)
		update.POST("/", h.AddMetricJSON)
		value.Use(middleware.SetContextPlain)
		value.GET("/:type/:name", h.GetMetric)
		value.POST("/", h.GetMetricJSON)
	}
	return router
}
