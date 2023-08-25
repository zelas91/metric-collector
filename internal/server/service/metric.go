package service

import (
	"errors"
	"fmt"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/types"
	"strconv"
	"strings"
)

var log = logger.New()

type MemService struct {
	repo repository.MemRepository
}

func NewMetricsService(repo repository.MemRepository) *MemService {
	return &MemService{repo: repo}
}

func (s *MemService) AddMetricsJSON(metric payload.Metrics) (*payload.Metrics, error) {
	switch strings.ToLower(metric.MType) {
	case types.CounterType:
		if metric.Delta == nil {
			return nil, errors.New("counter delta not found")
		}
		val := s.addMetricCounterJSON(metric.ID, *metric.Delta)
		return &payload.Metrics{ID: metric.ID, MType: metric.MType, Delta: &val}, nil
	case types.GaugeType:
		if metric.Value == nil {
			return nil, errors.New("counter delta not found")
		}
		val := s.addMetricGaugeSON(metric.ID, *metric.Value)
		return &payload.Metrics{ID: metric.ID, MType: metric.MType, Value: &val}, nil
	default:
		return nil, errors.New("type mem error")
	}
}
func (s *MemService) addMetricGaugeSON(name string, value float64) float64 {
	return s.repo.AddMetricGauge(name, value)
}

func (s *MemService) addMetricCounterJSON(name string, value int64) int64 {
	return s.repo.AddMetricCounter(name, value)
}

func (s *MemService) GetMetrics() (map[string]types.MetricTypeValue, error) {
	gauge, err := s.repo.GetByType(types.GaugeType)
	if err != nil {
		return nil, fmt.Errorf("internal server error. %v", err)
	}
	counter, err := s.repo.GetByType(types.CounterType)
	if err != nil {
		return nil, fmt.Errorf("internal server error. %v", err)
	}
	arraysMetric := make(map[string]types.MetricTypeValue, len(gauge)+len(counter))

	for key, value := range gauge {
		arraysMetric[key] = value
	}
	for key, value := range counter {
		arraysMetric[key] = value
	}
	return arraysMetric, nil
}

func (s *MemService) GetMetric(name, t string) (types.MetricTypeValue, error) {
	val := s.repo.ReadMetric(name, t)
	if val == nil {
		return nil, fmt.Errorf(" not found metrics  name=%s , type=%s , val = %v", name, t, val)
	}
	return val, nil
}

func (s *MemService) AddMetric(name, t string, value string) error {
	if ok := checkValid(t, value); !ok {
		return errors.New("not valid name or type ")
	}

	switch strings.ToLower(t) {
	case types.CounterType:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("convert string to int64 error=%v", err)
		}
		s.repo.AddMetricCounter(name, val)
	case types.GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("convert string to int64 error=%v", err)
		}
		s.repo.AddMetricGauge(name, val)
	}
	return nil
}

func checkValid(typ, value string) bool {
	if !isValue(value) || !isType(typ) {
		return false
	}
	return true
}
func isValue(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Debugf("not valid value=%s, error=%v", value, err)
	}
	return err == nil
}

func isType(mType string) bool {
	switch mType {

	case types.CounterType:
	case types.GaugeType:
	default:
		return false
	}
	return true

}
