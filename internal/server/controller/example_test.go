package controller_test

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/controller"
	mock_service "github.com/zelas91/metric-collector/internal/server/service/mocks"
	"log"
	"net/http"
)

func Example() {
	// create new handler
	h := controller.NewMetricHandler(&mock_service.MockService{})

	// configure routes
	router := gin.Default()
	router.POST("/json/update", h.AddMetricJSON)
	router.POST("/json/updates", h.AddMetrics)
	router.GET("/json/value", h.GetMetricJSON)

	router.GET("/simple/ping", h.Ping)
	router.GET("/simple/main", h.GetMetrics)
	router.POST("/simple/update", h.AddMetric)
	router.POST("/simple/value", h.GetMetric)

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal(err)
	}
}
