package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/metric-collector/internal/server/advicerrors"
	"github.com/zelas91/metric-collector/internal/server/utils/types"
	"testing"
)

func TestCheckUpdateMetric(t *testing.T) {
	tests := []struct {
		name   string
		method string
		arrays []string
		want   *advicerrors.AppError
	}{
		{
			name:   "test not found #1",
			method: "POST",
			arrays: []string{"update", "counter", "Random", "17"},
			want:   advicerrors.NewErrNotFound("not found"),
		},
		{
			name:   "test ok #2",
			method: "POST",
			arrays: []string{"host", "update", "counter", "Random", "17"},
		},
		{
			name:   "test bad request #3",
			method: "POST",
			arrays: []string{"host", "update", "counter", "Random", "none"},
			want:   advicerrors.NewErrBadRequest("bad request"),
		},
		{
			name:   "test allow method #4",
			method: "GET",
			arrays: []string{"host", "update", "counter", "Random", "none"},
			want:   advicerrors.NewErrMethodNotAllowed("method not allowed"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := CheckUpdateMetric(test.method, test.arrays)
			assert.IsType(t, test.want, res)
			if res != nil {
				assert.Equal(t, test.want, res)
			}
		})
	}
}

func Test_isType(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		strType string
	}{
		{
			name:    "test isType Gauge yes #1",
			want:    true,
			strType: types.GaugeType,
		},
		{
			name:    "test isType no #2",
			want:    false,
			strType: "Gauges",
		},
		{
			name:    "test isType Counter ok #3",
			want:    true,
			strType: types.CounterType,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, isType(test.strType))
		})
	}
}

func Test_isValue(t *testing.T) {
	tests := []struct {
		name  string
		want  bool
		value string
	}{
		{
			name:  "test float64 is value #1",
			want:  true,
			value: "12.5",
		}, {
			name:  "test int64 is value #2",
			want:  true,
			value: "12",
		},
		{
			name:  "test invalid is value #3",
			want:  false,
			value: "none",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, isValue(test.value))
		})
	}
}
