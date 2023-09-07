package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/zelas91/metric-collector/internal/server/types"
)

//go:generate mockgen -package mocks -destination=./mocks/mock_repository.go -source=repository.go -package=mock MemRepository
type MemRepository interface {
	AddMetricGauge(name string, value float64) float64
	AddMetricCounter(name string, value int64) int64
	ReadMetric(name string, t string) types.MetricTypeValue
	GetByType(t string) (map[string]types.MetricTypeValue, error)
	Ping() error
	Shutdown() error
}

type MemStorage struct {
	Gauge   map[string]types.MetricTypeValue `json:"gauge"`
	Counter map[string]types.MetricTypeValue `json:"counter"`
	db      *sql.DB
}

func (m *MemStorage) SetDB(db *sql.DB) {
	m.db = db
}

func NewMemStorage(db *sql.DB) *MemStorage {
	return &MemStorage{Gauge: make(map[string]types.MetricTypeValue),
		Counter: make(map[string]types.MetricTypeValue),
		db:      db,
	}
}

func (m *MemStorage) AddMetricGauge(name string, value float64) float64 {
	m.Gauge[name] = types.Gauge(value)
	return value
}

func (m *MemStorage) AddMetricCounter(name string, value int64) int64 {
	existingValue, ok := m.Counter[name]
	if ok {
		newValue := types.Counter(value) + (existingValue.(types.Counter))
		m.Counter[name] = newValue
	} else {
		m.Counter[name] = types.Counter(value)
	}

	return int64(m.Counter[name].(types.Counter))
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

// UnmarshalJSON implementation
func (m *MemStorage) UnmarshalJSON(bytes []byte) error { //
	mem := &struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}{}
	if err := json.Unmarshal(bytes, mem); err != nil {
		return err
	}
	m.Counter = make(map[string]types.MetricTypeValue)
	m.Gauge = make(map[string]types.MetricTypeValue)
	for key, value := range mem.Counter {
		m.Counter[key] = types.Counter(value)
	}
	for key, value := range mem.Gauge {
		m.Gauge[key] = types.Gauge(value)
	}
	return nil

}

func (m *MemStorage) Ping() error {
	return m.db.Ping()
}

func (m *MemStorage) Shutdown() error {
	return m.db.Close()
}
