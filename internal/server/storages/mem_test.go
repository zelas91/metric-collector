package storages

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zelas91/metric-collector/internal/server/types"
	"strconv"
	"testing"
)

func TestAddMetric(t *testing.T) {
	memStorage := &MemStorage{
		Gauge: map[string]types.Gauge{
			"cpu_usage":    types.Gauge(0.85),
			"memory_usage": types.Gauge(0.6),
		},
		Counter: map[string]types.Counter{
			"requests": types.Counter(100),
			"errors":   types.Counter(5),
		},
	}
	tests := []struct {
		name       string
		metricType string
		expected   string
	}{
		{
			name:       "Write Gauge metric",
			metricType: types.GaugeType,
			expected:   "12",
		},
		{
			name:       "Write Counter metric",
			metricType: types.CounterType,
			expected:   "13",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memStorage.AddMetric(test.name, test.metricType, test.expected)
			if test.metricType == types.GaugeType {
				val, err := strconv.ParseFloat(test.expected, 64)
				require.NoError(t, err)
				assert.Equal(t, types.Gauge(val), memStorage.Gauge[test.name])
			} else {
				val, err := strconv.ParseInt(test.expected, 10, 64)
				require.NoError(t, err)
				assert.Equal(t, types.Counter(val), memStorage.Counter[test.name])
			}
		})
	}
}

func TestReadMetric(t *testing.T) {
	memStorage := &MemStorage{
		Gauge: map[string]types.Gauge{
			"cpu_usage":    types.Gauge(0.85),
			"memory_usage": types.Gauge(0.6),
		},
		Counter: map[string]types.Counter{
			"requests": types.Counter(100),
			"errors":   types.Counter(5),
		},
	}

	testCases := []struct {
		name       string
		metricName string
		metricType string
		expected   types.MetricTypeValue
	}{
		{
			name:       "Read Gauge Metric",
			metricName: "cpu_usage",
			metricType: types.GaugeType,
			expected:   types.Gauge(0.85),
		},
		{
			name:       "Read Counter Metric",
			metricName: "requests",
			metricType: types.CounterType,
			expected:   types.Counter(100),
		},
		{
			name:       "Unknown Metric Type",
			metricName: "unknown_metric",
			metricType: "unknown",
			expected:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := memStorage.ReadMetric(tc.metricName, tc.metricType)
			assert.Equal(t, tc.expected, result)
		})
	}
}
