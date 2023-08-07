package storages

import (
	"github.com/sirupsen/logrus"
	"github.com/zelas91/metric-collector/internal/server/types"
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
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logrus.Debugf("convert string to int64 error=%v", err)
		}
		existingValue, ok := m.Counter[name]
		if ok {
			newValue := val + existingValue.Value
			m.Counter[name] = types.Counter{Value: newValue}
		} else {
			m.Counter[name] = types.Counter{Value: val}
		}
	case types.GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logrus.Debugf("convert string to float64 error=%v", err)
		}
		m.Gauge[name] = types.Gauge{Value: val}
	}
}

func (m *MemStorage) ReadMetric(name string, t string) interface{} {
	switch t {
	case types.GaugeType:
		val, ok := m.Gauge[name]
		if !ok {
			return nil
		}
		return val.Value
	case types.CounterType:
		val, ok := m.Counter[name]
		if !ok {
			return nil
		}
		return val.Value
	default:
		return nil
	}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{Gauge: make(map[string]types.Gauge),
		Counter: make(map[string]types.Counter),
	}
}
