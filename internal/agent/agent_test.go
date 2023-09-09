package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadStats(t *testing.T) {
	stats := NewStats()
	prevPollCount := stats.PollCount
	prevRandomValue := stats.RandomValue

	stats.ReadStats()

	assert.Equal(t, stats.PollCount, prevPollCount+1)
	assert.NotEqual(t, stats.RandomValue, prevRandomValue)

}
func TestGetGauges(t *testing.T) {
	stats := NewStats()
	gauges := stats.GetGauges()

	var expectedGauges = map[string]float64{
		"Alloc":         float64(stats.Alloc),
		"GCCPUFraction": float64(stats.GCCPUFraction),
		"GCSys":         float64(stats.GCSys),
		"RandomValue":   float64(stats.RandomValue),
	}
	for key, expectedValue := range expectedGauges {
		actualValue, ok := gauges[key]
		assert.True(t, ok, "Gauge with key %q not found", key)
		assert.Equal(t, expectedValue, actualValue)
	}
}
func TestGetCounters(t *testing.T) {
	stats := NewStats()
	counters := stats.GetCounters()

	expectedCounters := map[string]int64{
		"PollCount": stats.PollCount,
	}

	for key, expectedValue := range expectedCounters {
		actualValue, ok := counters[key]
		assert.True(t, ok, "Counter with key %q not found", key)
		assert.Equal(t, expectedValue, actualValue)
	}
}
