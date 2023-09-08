package repository

//go:generate mockgen -package mocks -destination=./mocks/mock_repository.go -source=repository.go -package=mock MemRepository
type MemRepository interface {
	AddMetricGauge(name string, value float64) float64
	AddMetricCounter(name string, value int64) int64
	GetMetricGauge(name string) *float64
	GetMetricCounter(name string) *int64
	GetMetricGauges() map[string]float64
	GetMetricCounters() map[string]int64
}

type MemStorage struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"` //name , type , value
}

func (m *MemStorage) AddMetricGauge(name string, value float64) float64 {
	m.Gauge[name] = value
	return value
}

func (m *MemStorage) AddMetricCounter(name string, value int64) int64 {
	existingValue, ok := m.Counter[name]
	if ok {
		newValue := value + existingValue
		m.Counter[name] = newValue
	} else {
		m.Counter[name] = value
	}

	return m.Counter[name]
}

func (m *MemStorage) GetMetricGauge(name string) *float64 {
	val, ok := m.Gauge[name]
	if !ok {
		return nil
	}
	return &val
}

func (m *MemStorage) GetMetricCounter(name string) *int64 {
	val, ok := m.Counter[name]
	if !ok {
		return nil
	}
	return &val
}

func NewMemStorage() *MemStorage {
	return &MemStorage{Gauge: make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *MemStorage) GetMetricGauges() map[string]float64 {
	return m.Gauge
}
func (m *MemStorage) GetMetricCounters() map[string]int64 {
	return m.Counter
}
