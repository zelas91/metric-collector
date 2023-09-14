package repository

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/metric-collector/internal/server/types"
	"testing"
)

var (
	cpuValue  = 20.07
	poolValue = int64(20)
)

func TestAddMetric(t *testing.T) {

	tests := []struct {
		name     string
		metric   Metric
		expected *Metric
	}{
		{
			name: "#1 Gauge OK",
			metric: Metric{
				ID:    "CPU",
				MType: types.GaugeType,
				Value: &cpuValue,
			},
			expected: &Metric{
				ID:    "CPU",
				MType: types.GaugeType,
				Value: &cpuValue,
			},
		}, {
			name: "#2 Counter OK",
			metric: Metric{
				ID:    "POOL",
				MType: types.CounterType,
				Delta: &poolValue,
			},
			expected: &Metric{
				ID:    "POOL",
				MType: types.CounterType,
				Delta: &poolValue,
			},
		}, {
			name: "#3 error type",
			metric: Metric{
				ID:    "CPUZ",
				Delta: &poolValue,
			},
			expected: &Metric{
				ID:    "CPUZ",
				Delta: &poolValue,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem := NewMemStorage()
			result := mem.AddMetric(context.Background(), test.metric)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetMetric(t *testing.T) {
	mem := MemStorage{
		mem: map[string]Metric{
			"CPU": {ID: "CPU",
				MType: types.GaugeType,
				Value: &cpuValue},
			"POOL": {ID: "POOL",
				MType: types.CounterType,
				Delta: &poolValue}},
	}
	type expected struct {
		metric *Metric
		err    error
	}
	tests := []struct {
		name     string
		metric   Metric
		expected expected
	}{
		{
			name: "#1 get cpu",
			metric: Metric{
				ID: "CPU",
			},
			expected: expected{
				metric: &Metric{
					ID:    "CPU",
					MType: types.GaugeType,
					Value: &cpuValue,
				},
				err: nil,
			},
		}, {
			name: "#2 get pool",
			metric: Metric{
				ID: "POOL",
			},
			expected: expected{
				metric: &Metric{
					ID:    "POOL",
					MType: types.CounterType,
					Delta: &poolValue,
				},
				err: nil,
			},
		}, {
			name: "#3 not found metrics",
			metric: Metric{
				ID: "P",
			},
			expected: expected{
				err: errors.New("not found metrics"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := mem.GetMetric(context.Background(), test.metric.ID)
			assert.Equal(t, test.expected.metric, result)
			assert.Equal(t, test.expected.err, err)
		})

	}
}
