package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/repository"
)

type Handler struct {
	MemStore repository.MemRepository
}

func NewHandler(memStore repository.MemRepository) *Handler {
	return &Handler{MemStore: memStore}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(MiddlewareSetContextHTML)
	router.GET("/", h.GetMetrics)
	update := router.Group("/update")
	value := router.Group("/value")
	{
		update.Use(MiddlewareSetContextPlain)
		update.POST("/:type/:name/:value", h.AddMetric)
		value.Use(MiddlewareSetContextPlain)
		value.GET("/:type/:name", h.GetMetric)
	}
	return router
}
