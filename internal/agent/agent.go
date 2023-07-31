package agent

import (
	"github.com/zelas91/metric-collector/internal/server/utils/types"
	"math/rand"
	"runtime"
)

type Stats struct {
	runtime.MemStats
	PollCount   int64
	RandomValue int
}

func NewStats() *Stats {
	return &Stats{PollCount: 0, RandomValue: 0}
}

func (s *Stats) ReadStats() {
	runtime.ReadMemStats(&s.MemStats)
	s.PollCount += 1
	s.RandomValue = rand.Int()
}
func (s *Stats) GetGauges() map[string]*types.Gauge {
	return map[string]*types.Gauge{
		"Alloc":         &types.Gauge{Value: float64(s.Alloc)},
		"GCCPUFraction": &types.Gauge{Value: s.GCCPUFraction},
		"GCSys":         &types.Gauge{Value: float64(s.GCSys)},
		"HeapAlloc":     &types.Gauge{Value: float64(s.HeapAlloc)},
		"BuckHashSys":   &types.Gauge{Value: float64(s.BuckHashSys)},
		"HeapIdle":      &types.Gauge{Value: float64(s.HeapIdle)},
		"HeapInuse":     &types.Gauge{Value: float64(s.HeapInuse)},
		"HeapObjects":   &types.Gauge{Value: float64(s.HeapObjects)},
		"HeapReleased":  &types.Gauge{Value: float64(s.HeapReleased)},
		"HeapSys":       &types.Gauge{Value: float64(s.HeapSys)},
		"LastGC":        &types.Gauge{Value: float64(s.LastGC)},
		"Lookups":       &types.Gauge{Value: float64(s.Lookups)},
		"MCacheInuse":   &types.Gauge{Value: float64(s.MCacheInuse)},
		"MCacheSys":     &types.Gauge{Value: float64(s.MCacheSys)},
		"MSpanInuse":    &types.Gauge{Value: float64(s.MSpanInuse)},
		"MSpanSys":      &types.Gauge{Value: float64(s.MSpanSys)},
		"Mallocs":       &types.Gauge{Value: float64(s.Mallocs)},
		"NextGC":        &types.Gauge{Value: float64(s.NextGC)},
		"NumForcedGC":   &types.Gauge{Value: float64(s.NumForcedGC)},
		"NumGC":         &types.Gauge{Value: float64(s.NumGC)},
		"OtherSys":      &types.Gauge{Value: float64(s.OtherSys)},
		"PauseTotalNs":  &types.Gauge{Value: float64(s.PauseTotalNs)},
		"StackInuse":    &types.Gauge{Value: float64(s.StackInuse)},
		"StackSys":      &types.Gauge{Value: float64(s.StackSys)},
		"Sys":           &types.Gauge{Value: float64(s.Sys)},
		"TotalAlloc":    &types.Gauge{Value: float64(s.TotalAlloc)},
		"RandomValue":   &types.Gauge{Value: float64(s.RandomValue)},
	}
}

func (s *Stats) GetCounters() map[string]*types.Counter {
	return map[string]*types.Counter{
		"PoolCounter": &types.Counter{Value: s.PollCount},
	}
}
