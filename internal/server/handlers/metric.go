package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/storages"
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
		"<html>       " +
		"<head>             " +
		"<title>Table Metrics</title>         " +
		"</head>         " +
		"<body>            " +
		"<table>               " +
		"<td>Name</td>                     " +
		"<td>Value</td> 	" +
		"{{range $key, $value := .}}                     " +
		"<tr>                         " +
		"<td>{{$key}}</td>                     " +
		"<td>{{$value}}</td>                    " +
		"</tr>                 " +
		"{{end}}             " +
		"</table>         " +
		"</body>         " +
		"</html>"
)

func (h *Handler) AddMetric(c *gin.Context) {
	val := c.Param(paramValue)
	t := c.Param(paramType)
	checkValid(c, t, val)
	h.MemStore.AddMetric(c.Param(paramName), t, val)
}

func (h *Handler) GetMetric(c *gin.Context) {
	t := c.Param(paramType)
	name := c.Param(paramName)
	result := h.MemStore.ReadMetric(name, t)
	if result == nil {
		newErrorResponse(c, http.StatusNotFound, "not found")
	}
	if _, err := c.Writer.WriteString(fmt.Sprintf("%v", result)); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) GetMetrics(c *gin.Context) {
	body, err := template.New("test").Parse(templateHTML)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	memStore, ok := h.MemStore.(*storages.MemStorage)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	arraysMetric := make(map[string]interface{})

	for key, value := range memStore.Gauge {
		arraysMetric[key] = value.Value
	}
	for key, value := range memStore.Counter {
		arraysMetric[key] = value.Value
	}

	if err = body.Execute(c.Writer, arraysMetric); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func checkValid(c *gin.Context, typ, value string) {
	if !isValue(value) || !isType(typ) {
		newErrorResponse(c, http.StatusBadRequest, "not valid name or type ")
		return
	}
}
func isValue(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
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
