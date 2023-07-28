package storages

import (
	"github.com/zelas91/metric-collector/internal/utils/types"
)

type MemStorage struct {
	metrics map[string]map[types.MetricType]interface{} //name , type , value
}

func (m *MemStorage) Metrics() map[string]map[types.MetricType]interface{} {
	return m.metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: make(map[string]map[types.MetricType]interface{})}
}
