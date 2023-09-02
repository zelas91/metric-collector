package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/types"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Service interface {
	AddMetric(name, t string, value string) error
	GetMetric(name, t string) types.MetricTypeValue
	GetMetrics() (map[string]types.MetricTypeValue, error)
	AddMetricsJSON(metric payload.Metrics) (*payload.Metrics, error)
}

var log = logger.New()
var once sync.Once

type MemService struct {
	repo repository.MemRepository
	cfg  *config.Config
	ctx  context.Context
}

func NewMetricsService(repo repository.MemRepository, cfg *config.Config, ctx context.Context) *MemService {
	mem := readMetricsDB(cfg)
	if mem != nil {
		return &MemService{repo: mem, cfg: cfg, ctx: ctx}
	}
	return &MemService{repo: repo, cfg: cfg, ctx: ctx}
}

func (s *MemService) AddMetricsJSON(metric payload.Metrics) (*payload.Metrics, error) {

	metrics := &payload.Metrics{
		ID:    metric.ID,
		MType: metric.MType,
	}

	switch strings.ToLower(metric.MType) {
	case types.CounterType:
		if metric.Delta == nil {
			return nil, errors.New("counter delta not found")
		}
		val := s.addMetricCounterJSON(metric.ID, *metric.Delta)
		metrics.Delta = &val
	case types.GaugeType:
		if metric.Value == nil {
			return nil, errors.New("gauge value not found")
		}
		val := s.addMetricGaugeJSON(metric.ID, *metric.Value)
		metrics.Value = &val
	default:
		return nil, errors.New("type mem error")
	}
	s.saveMetric()
	return metrics, nil
}

func (s *MemService) addMetricGaugeJSON(name string, value float64) float64 {
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

func (s *MemService) GetMetric(name, t string) types.MetricTypeValue {
	return s.repo.ReadMetric(name, t)
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

func readMetricsDB(cfg *config.Config) *repository.MemStorage {
	if cfg.Restore == nil || cfg.FilePath == nil {
		return nil
	}
	path := *cfg.FilePath

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	info, _ := file.Stat()
	data := make([]byte, info.Size())
	if err != nil {
		log.Errorf("open file err: %v", err)
		return nil
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("file close err : %v", err)
			return
		}
	}()

	if _, err := file.Read(data); err != nil {
		log.Errorf("read file err: %v", err)
		return nil
	}
	var mem *repository.MemStorage

	if err = json.Unmarshal(data, &mem); err != nil {
		log.Errorf("read metrics db err: %v", err)
		return nil
	}
	return mem
}

func (s *MemService) save() error {
	if s.cfg.FilePath == nil {
		return errors.New("file path = nil")
	}
	data, err := json.MarshalIndent(s.repo, "", " ")
	if err != nil {
		return err
	}
	path := *s.cfg.FilePath
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("file close err : %v", err)
			return
		}
	}()

	if _, err = file.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *MemService) asyncSave() {
	once.Do(func() {
		go func() {
			tickerStoreInterval := time.NewTicker(time.Duration(*s.cfg.StoreInterval) * time.Second)
			for {
				select {
				case <-tickerStoreInterval.C:
					if err := s.save(); err != nil {
						log.Errorf("save error %v", err)
					}
				case <-s.ctx.Done():
					return
				}

			}
		}()
	})
}

func (s *MemService) syncSave() {
	if err := s.save(); err != nil {
		log.Errorf("save error %v", err)
		return
	}
}
func (s *MemService) saveMetric() {
	if s.cfg.StoreInterval == nil {
		return
	}

	if *s.cfg.StoreInterval == 0 {
		s.syncSave()
		return
	}
	s.asyncSave()
}
