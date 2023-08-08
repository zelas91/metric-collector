package storages

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zelas91/metric-collector/internal/server/types"
	"strconv"
	"strings"
)

type MemStorage struct {
	Gauge   map[string]types.MetricTypeValue
	Counter map[string]types.MetricTypeValue //name , type , value
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
			newValue := types.Counter(val) + (existingValue.(types.Counter))
			m.Counter[name] = newValue
		} else {
			m.Counter[name] = types.Counter(val)
		}
	case types.GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logrus.Debugf("convert string to float64 error=%v", err)
		}
		m.Gauge[name] = types.Gauge(val)
	}
}

func (m *MemStorage) ReadMetric(name string, t string) types.MetricTypeValue {
	switch t {
	case types.GaugeType:
		val, ok := m.Gauge[name]
		if !ok {
			return nil
		}
		return val
	case types.CounterType:
		val, ok := m.Counter[name]
		if !ok {
			return nil
		}
		return val
	default:
		return nil
	}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{Gauge: make(map[string]types.MetricTypeValue),
		Counter: make(map[string]types.MetricTypeValue),
	}
}
func (m *MemStorage) GetByType(t string) (map[string]types.MetricTypeValue, error) {
	switch t {
	case types.GaugeType:
		return m.Gauge, nil
	case types.CounterType:
		return m.Counter, nil
	default:
		return nil, fmt.Errorf("type %s not found", t)

	}
}
