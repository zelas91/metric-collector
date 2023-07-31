package storages

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zelas91/metric-collector/internal/server/utils/types"
	"strconv"
	"testing"
)

func TestMemStorage_AddMetric(t *testing.T) {
	tests := []struct {
		name    string
		memType string
		value   string
		mem     *MemStorage
	}{
		{
			name:    "test Gauge 1",
			memType: types.GaugeType,
			value:   "12",
			mem:     NewMemStorage(),
		},
		{
			name:    "test Counter 2",
			memType: types.CounterType,
			value:   "13",
			mem:     NewMemStorage(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mem.AddMetric(test.name, test.memType, test.value)
			if test.memType == types.GaugeType {
				val, err := strconv.ParseFloat(test.value, 64)
				require.NoError(t, err)
				assert.Equal(t, val, test.mem.Gauge[test.name].Value)
			} else {
				val, err := strconv.ParseInt(test.value, 10, 64)
				require.NoError(t, err)
				assert.Equal(t, val, test.mem.Counter[test.name].Value)
			}
		})
	}
}
