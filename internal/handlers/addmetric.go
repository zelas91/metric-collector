package handlers

import (
	advUtils "github.com/zelas91/metric-collector/internal/advicerrors/utils"
	"github.com/zelas91/metric-collector/internal/storages"
	"github.com/zelas91/metric-collector/internal/utils/types"
	"net/http"
	"strconv"
	"strings"
)

type MetricHandler struct {
	mem *storages.MemStorage
}

func NewMetricHandler(mem *storages.MemStorage) *MetricHandler {
	return &MetricHandler{mem: mem}
}

//	func UpdTest(h *UpdateHandler) advicerrors.AppHandler {
//		return func(w http.ResponseWriter, r *http.Request) error {
//			parts := strings.Split(r.URL.Path, "/")
//			err := advUtils.CheckUpdateMetric(r.Method, parts)
//			if err != nil {
//				return err
//			}
//			metricTypeStr := parts[2]
//			name := parts[3]
//			value := parts[4]
//
//			metrics := h.mem.Metrics()
//			if _, ok := metrics[name]; !ok {
//				metrics[name] = make(map[types.MetricType]interface{})
//			}
//			switch strings.ToLower(metricTypeStr) {
//			case "counter":
//				val, _ := strconv.ParseInt(value, 10, 64)
//				existingValue, ok := metrics[name][types.Counter]
//				if ok {
//					newValue := val + existingValue.(int64)
//					metrics[name][types.Counter] = newValue
//				} else {
//					metrics[name][types.Counter] = val
//				}
//			case "gauge":
//				val, _ := strconv.ParseFloat(value, 64)
//				metrics[name][types.Gauge] = val
//			}
//			return nil
//		}
//	}
//
// post method
func (h *MetricHandler) MetricAdd(w http.ResponseWriter, r *http.Request) error {
	parts := strings.Split(r.URL.Path, "/")
	err := advUtils.CheckUpdateMetric(r.Method, parts)
	if err != nil {
		return err
	}
	metricTypeStr := parts[2]
	name := parts[3]
	value := parts[4]

	metrics := h.mem.Metrics()
	if _, ok := metrics[name]; !ok {
		metrics[name] = make(map[types.MetricType]interface{})
	}
	switch strings.ToLower(metricTypeStr) {
	case "counter":
		val, _ := strconv.ParseInt(value, 10, 64)
		existingValue, ok := metrics[name][types.Counter]
		if ok {
			newValue := val + existingValue.(int64)
			metrics[name][types.Counter] = newValue
		} else {
			metrics[name][types.Counter] = val
		}
	case "gauge":
		val, _ := strconv.ParseFloat(value, 64)
		metrics[name][types.Gauge] = val
	}
	return nil
}
