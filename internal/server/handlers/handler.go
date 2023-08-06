package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/middleware"
)

func InitRoutes(h *controller.MetricHandler) *gin.Engine {
	router := gin.New()

	router.Use(middleware.SetContextHTML)
	router.GET("/", h.GetMetrics)
	update := router.Group("/update")
	value := router.Group("/value")
	{
		update.Use(middleware.SetContextPlain)
		update.POST("/:type/:name/:value", h.AddMetric)
		value.Use(middleware.SetContextPlain)
		value.GET("/:type/:name", h.GetMetric)
	}
	return router
}
