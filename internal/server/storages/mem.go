package storages

import (
	"encoding/json"
	"fmt"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/types"
	"os"
	"sync"
	"time"
)

var (
	cfg  *config.Config
	log  = logger.New()
	once sync.Once
)

type MemStorage struct {
	Gauge   map[string]types.MetricTypeValue `json:"gauge"`
	Counter map[string]types.MetricTypeValue `json:"counter"` //name , type , value
}

func NewMemStorage(config *config.Config) *MemStorage {
	cfg = config
	if *cfg.Restore {
		if mem := readMetricsDB(); mem != nil {
			return mem
		}
	}

	return &MemStorage{Gauge: make(map[string]types.MetricTypeValue),
		Counter: make(map[string]types.MetricTypeValue),
	}
}

func (m *MemStorage) UnmarshalJSON(bytes []byte) error {
	type memStorageAlias MemStorage
	mem := &struct {
		*memStorageAlias
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}{memStorageAlias: (*memStorageAlias)(m)}
	if err := json.Unmarshal(bytes, mem); err != nil {
		return err
	}
	m.Counter = make(map[string]types.MetricTypeValue)
	m.Gauge = make(map[string]types.MetricTypeValue)
	for key, value := range mem.Counter {
		m.Counter[key] = types.Counter(value)
	}
	for key, value := range mem.Gauge {
		m.Gauge[key] = types.Gauge(value)
	}
	return nil

}

func (m *MemStorage) AddMetricGauge(name string, value float64) float64 {
	m.Gauge[name] = types.Gauge(value)
	return value
}

func (m *MemStorage) AddMetricCounter(name string, value int64) int64 {
	existingValue, ok := m.Counter[name]
	if ok {
		newValue := types.Counter(value) + (existingValue.(types.Counter))
		m.Counter[name] = newValue
	} else {
		m.Counter[name] = types.Counter(value)
	}

	return int64(m.Counter[name].(types.Counter))
}

func (m *MemStorage) ReadMetric(name string, t string) types.MetricTypeValue {
	switch t {
	case types.GaugeType:
		val, ok := m.Gauge[name]
		if !ok {
			return nil
		}
		return val
	case types.CounterType:
		val, ok := m.Counter[name]
		if !ok {
			return nil
		}
		return val
	default:
		return nil
	}
}

func (m *MemStorage) save() error {
	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile("metrics-db.json", data, 0666)
}
func (m *MemStorage) saveMetric() {
	if *cfg.StoreInterval == 0 {
		log.Debug("start sync")
		if err := m.save(); err != nil {
			log.Debugf("SAVE ERROR %v", err)
		}

		return
	}

	once.Do(func() {
		log.Debugf("START interval")
		go func() {
			tickerStoreInterval := time.NewTicker(time.Duration(*cfg.StoreInterval) * time.Second)
			for {
				<-tickerStoreInterval.C
				log.Debug("SAVEEEEEEEE")
				if err := m.save(); err != nil {
					log.Debugf("SAVE ERROR %v", err)
				}
			}
		}()
	})
}

func readMetricsDB() *MemStorage {
	data, err := os.ReadFile("metrics-db.json")
	if err != nil {
		log.Debugf("read file err: %v", err)
		return nil
	}

	var mem *MemStorage

	err = json.Unmarshal(data, &mem)
	if err != nil {
		log.Debugf("ERROR UNMARSHAL %v", err)
		return nil
	}
	return mem
}

func (m *MemStorage) GetByType(t string) (map[string]types.MetricTypeValue, error) {
	switch t {
	case types.GaugeType:
		return m.Gauge, nil
	case types.CounterType:
		return m.Counter, nil
	default:
		return nil, fmt.Errorf("type %s not found", t)

	}
}
