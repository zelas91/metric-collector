// Package repository working with data in memory, file or database format
package repository

import "context"

// StorageRepository - interface for working with the metrics repository.
//
//go:generate mockgen -package mocks -destination=./mocks/mock_repository.go -source=repository.go -package=mock StorageRepository
type StorageRepository interface {
	// AddMetric adds a metric to the storage.
	AddMetric(ctx context.Context, metrics Metric) *Metric
	// GetMetric gets a metric from the repository by name.
	GetMetric(ctx context.Context, name string) (*Metric, error)
	// GetMetrics gets all metrics from the repository.
	GetMetrics(ctx context.Context) []Metric
	// AddMetrics adds a metric slice to the repository.
	AddMetrics(ctx context.Context, metrics []Metric) error
}

// Shutdown - interface for shutdown.
type Shutdown interface {
	// Shutdown performs shutdown.
	Shutdown() error
}

// Ping - interface for checking the availability of the service.
type Ping interface {
	// Ping performs a service availability check.
	Ping() error
}

// Metric - the structure representing the metric.
type Metric struct {
	ID    string   `json:"id" db:"name"`                         // имя метрики
	MType string   `json:"type" db:"type"`                       // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" db:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value,omitempty"` // значение метрики в случае передачи gauge
}
