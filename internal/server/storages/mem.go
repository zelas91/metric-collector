package storages

import (
	"github.com/zelas91/metric-collector/internal/server/utils/types"
	"strconv"
	"strings"
)

type MemStorage struct {
	Gauge   map[string]types.Gauge
	Counter map[string]types.Counter //name , type , value
}

func (m *MemStorage) AddMetric(name, typeMetric, value string) {
	switch strings.ToLower(typeMetric) {
	case types.CounterType:
		val, _ := strconv.ParseInt(value, 10, 64)
		existingValue, ok := m.Counter[name]
		if ok {
			newValue := val + existingValue.Value
			m.Counter[name] = types.Counter{Value: newValue}
		} else {
			m.Counter[name] = types.Counter{Value: val}
		}
	case types.GaugeType:
		val, _ := strconv.ParseFloat(value, 64)
		m.Gauge[name] = types.Gauge{Value: val}
	}
}

func (m *MemStorage) ReadMetric(name string) map[string]interface{} {
	//TODO implement me
	panic("implement me")
}

func NewMemStorage() *MemStorage {
	return &MemStorage{Gauge: make(map[string]types.Gauge),
		Counter: make(map[string]types.Counter),
	}
}
