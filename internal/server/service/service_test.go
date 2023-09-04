package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/repository"
	mock "github.com/zelas91/metric-collector/internal/server/repository/mocks"
	"github.com/zelas91/metric-collector/internal/server/types"
	"testing"
)

func TestAddMetricJSON(t *testing.T) {
	type result struct {
		excepted payload.Metrics
		err      error
	}
	serv := NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())
	gaugeValue := 20.123
	deltaValue := int64(200)
	tests := []struct {
		name  string
		want  result
		sense payload.Metrics
	}{
		{
			name: "# 1 Gauge Add JSON",
			want: result{
				excepted: payload.Metrics{
					ID:    "CPU",
					MType: types.GaugeType,
					Value: &gaugeValue,
				},
				err: nil,
			},
			sense: payload.Metrics{
				ID:    "CPU",
				MType: types.GaugeType,
				Value: &gaugeValue,
			},
		},
		{
			name: "# 2 Counter Add JSON",
			want: result{
				excepted: payload.Metrics{
					ID:    "Poll",
					MType: types.CounterType,
					Delta: &deltaValue,
				},
				err: nil,
			},
			sense: payload.Metrics{
				ID:    "Poll",
				MType: types.CounterType,
				Delta: &deltaValue,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := serv.AddMetricsJSON(test.sense)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, &test.want.excepted, res)
		})
	}
}
func TestIsType(t *testing.T) {
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

func TestIsValue(t *testing.T) {
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

func TestAddMetricGaugeMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memStorage := mock.NewMockMemRepository(ctrl)
	memStorage.EXPECT().AddMetricCounter("test", int64(20)).Return(int64(20))
	memStorage.EXPECT().GetByType("gauge").Return(map[string]types.MetricTypeValue{"test3": types.Gauge(15.7)}, nil)
	memStorage.EXPECT().GetByType("counter").Return(map[string]types.MetricTypeValue{"test4": types.Counter(15)}, nil)
	serv := NewMetricsService(memStorage, &config.Config{}, context.Background())
	fmt.Println(serv.GetMetrics())
	fmt.Println(serv.AddMetric("test", "counter", "20"))
}
