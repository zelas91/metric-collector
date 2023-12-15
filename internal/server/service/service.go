// Package service use case
package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/types"
)

//go:generate mockgen -package mocks -destination=./mocks/mock_service.go -source=service.go -package=mock_service Service
type Service interface {
	AddMetric(ctx context.Context, name, mType, value string) (*repository.Metric, error)
	GetMetric(ctx context.Context, name string) (*repository.Metric, error)
	GetMetrics(ctx context.Context) []repository.Metric
	AddMetricJSON(ctx context.Context, metric repository.Metric) (*repository.Metric, error)
	AddMetrics(ctx context.Context, metrics []repository.Metric) error
}

var (
	log = logger.New()
)

type MemService struct {
	repo repository.StorageRepository
	cfg  *config.Config
	ctx  context.Context
}

func NewMemService(ctx context.Context, repo repository.StorageRepository, cfg *config.Config) *MemService {
	return &MemService{repo: repo, cfg: cfg, ctx: ctx}
}

func (m *MemService) AddMetric(ctx context.Context, name, mType, value string) (*repository.Metric, error) {
	if !checkValid(mType, value) {
		return nil, errors.New("not valid name or type ")
	}
	metric := repository.Metric{ID: name, MType: mType}
	switch strings.ToLower(mType) {
	case types.CounterType:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("service addmetric type=%s : error=%w", mType, err)
		}
		metric.Delta = &val
	case types.GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("service addmetric type=%s : error=%w", mType, err)
		}
		metric.Value = &val
	}
	return m.repo.AddMetric(ctx, metric), nil
}

func (m *MemService) GetMetric(ctx context.Context, name string) (*repository.Metric, error) {
	return m.repo.GetMetric(ctx, name)
}

func (m *MemService) GetMetrics(ctx context.Context) []repository.Metric {
	return m.repo.GetMetrics(ctx)
}

func (m *MemService) AddMetricJSON(ctx context.Context, metric repository.Metric) (*repository.Metric, error) {
	switch metric.MType {
	case types.GaugeType:
		if metric.Value == nil {
			return nil, errors.New("gauge value not found")
		}
		return m.repo.AddMetric(ctx, metric), nil
	case types.CounterType:
		if metric.Delta == nil {
			return nil, errors.New("counter delta not found")
		}
		return m.repo.AddMetric(ctx, metric), nil
	default:
		return nil, errors.New("type mem error")
	}
}

func (m *MemService) Ping() error {
	repo, ok := m.repo.(repository.Ping)
	if ok {
		return repo.Ping()
	}
	return errors.New("not implementation  interface Ping")
}

func (m *MemService) AddMetrics(ctx context.Context, metrics []repository.Metric) error {
	return m.repo.AddMetrics(ctx, metrics)
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
		log.Errorf("not valid value=%s, error=%v", value, err)
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
