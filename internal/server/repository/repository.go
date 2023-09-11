package repository

//go:generate mockgen -package mocks -destination=./mocks/mock_Repository.go -source=repository.go -package=mock StorageRepository
type StorageRepository interface {
	AddMetric(metrics Metric) *Metric
	GetMetric(name string) (*Metric, error)
	GetMetrics() []Metric
}

type Shutdown interface {
	Shutdown() error
}
type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
