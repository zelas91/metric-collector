package repository

import "github.com/zelas91/metric-collector/internal/server/types"

type MemRepository interface {
	AddMetricGauge(name string, value float64) float64
	AddMetricCounter(name string, value int64) int64
	ReadMetric(name string, t string) types.MetricTypeValue
	GetByType(t string) (map[string]types.MetricTypeValue, error)
}
