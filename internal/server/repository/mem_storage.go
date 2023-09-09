package repository

import (
	"errors"
	"github.com/zelas91/metric-collector/internal/server/types"
)

type MemStore struct {
	mem map[string]Metric
}

func NewMemStore() *MemStore {
	return &MemStore{
		mem: make(map[string]Metric),
	}
}

func (m *MemStore) AddMetric(metric Metric) *Metric {
	switch metric.MType {
	case types.GaugeType:
		m.mem[metric.ID] = metric
	case types.CounterType:
		existingValue, ok := m.mem[metric.ID]
		if ok {
			newValue := *metric.Delta + *existingValue.Delta
			metric.Delta = &newValue
			m.mem[metric.ID] = metric
		} else {
			m.mem[metric.ID] = metric
		}
	}
	return &metric
}

func (m *MemStore) GetMetric(name string) (*Metric, error) {
	metric, ok := m.mem[name]
	if ok {
		return &metric, nil
	}
	return nil, errors.New("not found metrics")
}

func (m *MemStore) GetMetrics() []Metric {
	metrics := make([]Metric, len(m.mem))
	for _, val := range m.mem {
		metrics = append(metrics, val)
	}
	return metrics
}
