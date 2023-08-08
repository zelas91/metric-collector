package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/types"
	"html/template"
	"net/http"
	"strconv"
)

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
	MemStore repository.MemRepository
}

func NewMetricHandler(memStore repository.MemRepository) *MetricHandler {
	return &MetricHandler{MemStore: memStore}
}
func (h *MetricHandler) AddMetric(c *gin.Context) {
	val := c.Param(paramValue)
	t := c.Param(paramType)
	if ok := checkValid(t, val); !ok {
		payload.NewErrorResponse(c, http.StatusBadRequest, "not valid name or type ")
		return
	}
	h.MemStore.AddMetric(c.Param(paramName), t, val)
}

func (h *MetricHandler) GetMetric(c *gin.Context) {
	t := c.Param(paramType)
	name := c.Param(paramName)
	result := h.MemStore.ReadMetric(name, t)
	if result == nil {
		payload.NewErrorResponse(c, http.StatusNotFound, "not found")
		return
	}
	if _, err := c.Writer.WriteString(fmt.Sprintf("%v", result)); err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
}

func (h *MetricHandler) GetMetrics(c *gin.Context) {
	body, err := template.New("test").Parse(templateHTML)
	if err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	gauge, err := h.MemStore.GetByType(types.GaugeType)
	if err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	counter, err := h.MemStore.GetByType(types.CounterType)

	if err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	arraysMetric := make(map[string]types.MetricTypeValue, len(gauge)+len(counter))

	for key, value := range gauge {
		arraysMetric[key] = value
	}
	for key, value := range counter {
		arraysMetric[key] = value
	}

	if err = body.Execute(c.Writer, arraysMetric); err != nil {
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func checkValid(typ, value string) bool {
	if !isValue(value) || !isType(typ) {
		return false
	}
	return true
}
func isValue(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logrus.Debugf("not valid value=%s, error=%v", value, err)
	}
	return err == nil
}

func isType(mType string) bool {
	switch mType {

	case types.CounterType:
	case types.GaugeType:
	default:
		return false
	}
	return true

}
