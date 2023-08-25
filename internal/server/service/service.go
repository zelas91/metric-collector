package service

import (
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/types"
)

type Service interface {
	AddMetric(name, t string, value string) error
	GetMetric(name, t string) (types.MetricTypeValue, error)
	GetMetrics() (map[string]types.MetricTypeValue, error)
	AddMetricsJSON(metric payload.Metrics) (*payload.Metrics, error)
}
