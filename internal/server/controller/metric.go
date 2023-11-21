package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/server/types"
)

var log = logger.New()

const (
	paramName  = "name"
	paramType  = "type"
	paramValue = "value"
)
const (
	templateHTML = "<!DOCTYPE html> " +
		`<html>       
			<head>             
				<title>Table Metrics</title>         
			</head>         
			<body>            
				<table>               
					<td>Name</td>                     
					<td>Value</td> 	
					{{range $key, $value := .}}                     
					<tr>                         
						<td>{{$key}}</td>                    
						<td>{{$value}}</td>                    
					</tr>                 
					{{end}}             
				</table>         
			</body>         
		</html>`
)

type MetricHandler struct {
	memService service.Service
}

func NewMetricHandler(memService service.Service) *MetricHandler {
	return &MetricHandler{memService: memService}
}
func (h *MetricHandler) AddMetric(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	value := c.Param(paramValue)
	t := c.Param(paramType)

	if _, err := h.memService.AddMetric(c.Request.Context(), c.Param(paramName), t, value); err != nil {
		payload.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *MetricHandler) GetMetric(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	t := c.Param(paramType)
	name := c.Param(paramName)
	mem, err := h.memService.GetMetric(c.Request.Context(), name)
	if err != nil || !strings.EqualFold(t, mem.MType) {
		payload.NewErrorResponse(c, http.StatusNotFound, "not found")
		return
	}
	var result string
	switch t {
	case types.GaugeType:
		result = fmt.Sprintf("%v", *mem.Value)
	case types.CounterType:
		result = fmt.Sprintf("%v", *mem.Delta)
	}

	if _, err := c.Writer.WriteString(result); err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
}

func (h *MetricHandler) GetMetrics(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	body, err := template.New("values").Parse(templateHTML)
	if err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	metrics := h.memService.GetMetrics(c.Request.Context())
	mapMetrics := make(map[string]interface{}, len(metrics))
	for _, metric := range metrics {
		switch metric.MType {
		case types.GaugeType:
			mapMetrics[metric.ID] = *metric.Value
		case types.CounterType:
			mapMetrics[metric.ID] = *metric.Delta
		}
	}

	if err = body.Execute(c.Writer, mapMetrics); err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *MetricHandler) GetMetricJSON(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/json" {
		payload.NewErrorResponseJSON(c, http.StatusUnsupportedMediaType, "incorrect media type ")
		return
	}
	var request repository.Metric
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("request json error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	val, err := h.memService.GetMetric(c.Request.Context(), request.ID)
	if err != nil {
		log.Errorf("controller get metric error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusNotFound, err.Error())
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, val)
}
func (h *MetricHandler) AddMetricJSON(c *gin.Context) {

	if c.GetHeader("Content-Type") != "application/json" {
		payload.NewErrorResponseJSON(c, http.StatusUnsupportedMediaType, "incorrect media type ")
		return
	}

	var request repository.Metric
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("request json  error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.memService.AddMetricJSON(c.Request.Context(), request)
	if err != nil {
		log.Errorf("add metric json error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (h *MetricHandler) Ping(c *gin.Context) {
	ser, ok := h.memService.(repository.Ping)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := ser.Ping(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (h *MetricHandler) AddMetrics(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/json" {
		payload.NewErrorResponseJSON(c, http.StatusUnsupportedMediaType, "incorrect media type ")
		return
	}
	var request []repository.Metric
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("request json  error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.memService.AddMetrics(c.Request.Context(), request); err != nil {
		log.Errorf("add metrics err: %v", err)
		payload.NewErrorResponseJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

}
