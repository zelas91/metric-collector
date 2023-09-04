package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/repository"
	mock "github.com/zelas91/metric-collector/internal/server/repository/mocks"
	"github.com/zelas91/metric-collector/internal/server/types"
	"strconv"
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
	type mockBehavior func(s *mock.MockMemRepository, name, t, value string)
	type mem struct {
		name  string
		t     string
		value string
	}
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		mem          mem
		want         error
	}{
		{
			name: "# OK",
			mockBehavior: func(s *mock.MockMemRepository, name, t, value string) {
				switch t {
				case types.GaugeType:
					val, err := strconv.ParseFloat(value, 64)
					if err != nil {
						log.Fatalf("parsing float err : %v ", err)
					}
					s.EXPECT().AddMetricGauge(name, val).Return(val)
				case types.CounterType:
					val, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						log.Fatalf("parsing int64 err : %v ", err)
					}
					s.EXPECT().AddMetricCounter(name, val).Return(val)
				}

			},
			mem: mem{
				name:  "testCounter",
				value: "20",
				t:     types.CounterType,
			},
			want: nil,
		}, {
			name: "# not valid type",
			mockBehavior: func(s *mock.MockMemRepository, name, t, value string) {

			},
			mem: mem{
				name:  "testCounter",
				value: "20.0",
				t:     types.CounterType,
			},
			want: errors.New("convert string to int64 error=strconv.ParseInt: parsing \"20.0\": invalid syntax"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock.NewMockMemRepository(ctrl)
			test.mockBehavior(repo, test.mem.name, test.mem.t, test.mem.value)
			service := NewMetricsService(repo, &config.Config{}, context.Background())

			err := service.AddMetric(test.mem.name, test.mem.t, test.mem.value)
			assert.Equal(t, test.want, err)
		})
	}
}
