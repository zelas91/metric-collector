package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/repository"
	mock "github.com/zelas91/metric-collector/internal/server/repository/mocks"
	"github.com/zelas91/metric-collector/internal/server/types"
	"strconv"
	"testing"
)

func TestAddMetricJSON(t *testing.T) {
	type result struct {
		excepted repository.Metric
		err      error
	}
	serv := NewMemService(context.Background(), repository.NewMemStore(), &config.Config{})
	gaugeValue := 20.123
	deltaValue := int64(200)
	tests := []struct {
		name  string
		want  result
		sense repository.Metric
	}{
		{
			name: "# 1 Gauge Add JSON",
			want: result{
				excepted: repository.Metric{
					ID:    "CPU",
					MType: types.GaugeType,
					Value: &gaugeValue,
				},
				err: nil,
			},
			sense: repository.Metric{
				ID:    "CPU",
				MType: types.GaugeType,
				Value: &gaugeValue,
			},
		},
		{
			name: "# 2 Counter Add JSON",
			want: result{
				excepted: repository.Metric{
					ID:    "Poll",
					MType: types.CounterType,
					Delta: &deltaValue,
				},
				err: nil,
			},
			sense: repository.Metric{
				ID:    "Poll",
				MType: types.CounterType,
				Delta: &deltaValue,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := serv.AddMetricJSON(test.sense)

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

func TestAddMetric(t *testing.T) {
	deltaValue1 := int64(20)
	type mockBehavior func(s *mock.MockStorageRepository, name, t, value string)
	type mem struct {
		name  string
		t     string
		value string
	}
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		mem          mem
		wantErr      error
		expected     *repository.Metric
	}{
		{
			name: "# OK",
			mockBehavior: func(s *mock.MockStorageRepository, name, t, value string) {
				metric := repository.Metric{ID: name, MType: t}
				switch t {
				case types.GaugeType:
					val, err := strconv.ParseFloat(value, 64)
					if err != nil {
						log.Errorf("error convert string to float64, err:%v", err)
					}
					metric.Value = &val
					s.EXPECT().AddMetric(metric).Return(&metric)
				case types.CounterType:
					val, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						log.Errorf("error convert string to float64, err:%v", err)
					}
					metric.Delta = &val
					s.EXPECT().AddMetric(metric).Return(&metric)
				}

			},
			mem: mem{
				name:  "testCounter",
				value: "20",
				t:     types.CounterType,
			},
			expected: &repository.Metric{
				ID:    "testCounter",
				Delta: &deltaValue1,
				MType: types.CounterType,
			},
			wantErr: nil,
		}, {
			name: "# not valid type",
			mockBehavior: func(s *mock.MockStorageRepository, name, t, value string) {

			},
			mem: mem{
				name:  "testCounter",
				value: "20.0",
				t:     types.CounterType,
			},
			wantErr: errors.New("service addmetric type=counter : error=strconv.ParseInt: parsing \"20.0\": invalid syntax"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock.NewMockStorageRepository(ctrl)
			test.mockBehavior(repo, test.mem.name, test.mem.t, test.mem.value)
			service := NewMemService(context.Background(), repo, &config.Config{})

			metric, err := service.AddMetric(test.mem.name, test.mem.t, test.mem.value)
			if test.wantErr != nil {
				assert.Equal(t, test.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, metric)
		})
	}
}
