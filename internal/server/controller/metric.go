package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/server/types"
	"html/template"
	"net/http"
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
	MemService service.Service
}

func NewMetricHandler(MemService service.Service) *MetricHandler {
	return &MetricHandler{MemService: MemService}
}
func (h *MetricHandler) AddMetric(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	value := c.Param(paramValue)
	t := c.Param(paramType)
	if err := h.MemService.AddMetric(c.Param(paramName), t, value); err != nil {
		payload.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *MetricHandler) GetMetric(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	t := c.Param(paramType)
	name := c.Param(paramName)
	result, err := h.MemService.GetMetric(name, t)
	if err != nil {
		payload.NewErrorResponse(c, http.StatusNotFound, "not found")
		return
	}
	if _, err := c.Writer.WriteString(fmt.Sprintf("%v", result)); err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
}

func (h *MetricHandler) GetMetrics(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	body, err := template.New("test").Parse(templateHTML)
	if err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	arraysMetric, err := h.MemService.GetMetrics()
	if err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if err = body.Execute(c.Writer, arraysMetric); err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *MetricHandler) GetMetricJSON(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/json" {
		payload.NewErrorResponseJSON(c, http.StatusUnsupportedMediaType, "incorrect media type ")
		return
	}
	var request payload.Metrics
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Debugf("bind json  json error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	val, err := h.MemService.GetMetric(request.ID, request.MType)
	if err != nil {
		log.Debugf("get metric=%v error=%v ", request, err)
		payload.NewErrorResponseJSON(c, http.StatusNotFound, err.Error())
		return
	}
	result := payload.Metrics{
		ID:    request.ID,
		MType: request.MType,
	}
	switch request.MType {
	case types.CounterType:
		delta := int64(val.(types.Counter))
		result.Delta = &delta
	case types.GaugeType:
		value := float64(val.(types.Gauge))
		result.Value = &value
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func (h *MetricHandler) AddMetricJSON(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/json" {
		payload.NewErrorResponseJSON(c, http.StatusUnsupportedMediaType, "incorrect media type ")
		return
	}
	var request payload.Metrics
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Debugf("bind json  error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.MemService.AddMetricsJSON(request)
	if err != nil {
		log.Debugf("add metric json error=%v ", err)
		payload.NewErrorResponseJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}
