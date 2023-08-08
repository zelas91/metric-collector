package repository

import "github.com/zelas91/metric-collector/internal/server/types"

type MemRepository interface {
	AddMetric(name, typeMetric, value string)
	ReadMetric(name string, t string) types.MetricTypeValue
	GetByType(t string) (map[string]types.MetricTypeValue, error)
}
