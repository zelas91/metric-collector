package repository

//
//import (
//	"github.com/stretchr/testify/assert"
//	"github.com/zelas91/metric-collector/internal/server/types"
//	"testing"
//)
//
//func TestAddMetricGauge(t *testing.T) {
//	memStorage := &MemStorage{
//		Gauge: map[string]types.MetricTypeValue{
//			"cpu_usage":    types.Gauge(0.85),
//			"memory_usage": types.Gauge(0.6),
//		},
//		Counter: map[string]types.MetricTypeValue{
//			"requests": types.Counter(100),
//			"errors":   types.Counter(5),
//		},
//	}
//	tests := []struct {
//		name       string
//		metricType string
//		expected   float64
//	}{
//		{
//			name:       "Write Gauge metric",
//			metricType: types.GaugeType,
//			expected:   12.8,
//		},
//		{
//			name:       "Write Counter metric",
//			metricType: types.CounterType,
//			expected:   13,
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			memStorage.AddMetricGauge(test.name, test.expected)
//			assert.Equal(t, types.Gauge(test.expected), memStorage.Gauge[test.name])
//
//		})
//	}
//}
//
//func TestReadMetric(t *testing.T) {
//	memStorage := &MemStorage{
//		Gauge: map[string]types.MetricTypeValue{
//			"cpu_usage":    types.Gauge(0.85),
//			"memory_usage": types.Gauge(0.6),
//		},
//		Counter: map[string]types.MetricTypeValue{
//			"requests": types.Counter(100),
//			"errors":   types.Counter(5),
//		},
//	}
//
//	testCases := []struct {
//		name       string
//		metricName string
//		metricType string
//		expected   types.MetricTypeValue
//	}{
//		{
//			name:       "Read Gauge Metric",
//			metricName: "cpu_usage",
//			metricType: types.GaugeType,
//			expected:   types.Gauge(0.85),
//		},
//		{
//			name:       "Read Counter Metric",
//			metricName: "requests",
//			metricType: types.CounterType,
//			expected:   types.Counter(100),
//		},
//		{
//			name:       "Unknown Metric Type",
//			metricName: "unknown_metric",
//			metricType: "unknown",
//			expected:   nil,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			result := memStorage.ReadMetric(tc.metricName, tc.metricType)
//			assert.Equal(t, tc.expected, result)
//		})
//	}
//}
