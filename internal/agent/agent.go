package agent

import (
	"github.com/zelas91/metric-collector/internal/server/types"
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
func (s *Stats) GetGauges() map[string]types.Gauge {
	return map[string]types.Gauge{
		"Alloc":         types.Gauge(s.Alloc),
		"GCSys":         types.Gauge(s.GCSys),
		"HeapAlloc":     types.Gauge(s.HeapAlloc),
		"BuckHashSys":   types.Gauge(s.BuckHashSys),
		"GCCPUFraction": types.Gauge(s.GCCPUFraction),
		"HeapIdle":      types.Gauge(s.HeapIdle),
		"HeapInuse":     types.Gauge(s.HeapInuse),
		"HeapObjects":   types.Gauge(s.HeapObjects),
		"HeapReleased":  types.Gauge(s.HeapReleased),
		"HeapSys":       types.Gauge(s.HeapSys),
		"LastGC":        types.Gauge(s.LastGC),
		"Lookups":       types.Gauge(s.Lookups),
		"MCacheInuse":   types.Gauge(s.MCacheInuse),
		"MCacheSys":     types.Gauge(s.MCacheSys),
		"MSpanInuse":    types.Gauge(s.MSpanInuse),
		"MSpanSys":      types.Gauge(s.MSpanSys),
		"Mallocs":       types.Gauge(s.Mallocs),
		"NextGC":        types.Gauge(s.NextGC),
		"NumForcedGC":   types.Gauge(s.NumForcedGC),
		"NumGC":         types.Gauge(s.NumGC),
		"OtherSys":      types.Gauge(s.OtherSys),
		"PauseTotalNs":  types.Gauge(s.PauseTotalNs),
		"StackInuse":    types.Gauge(s.StackInuse),
		"StackSys":      types.Gauge(s.StackSys),
		"Sys":           types.Gauge(s.Sys),
		"TotalAlloc":    types.Gauge(s.TotalAlloc),
		"RandomValue":   types.Gauge(s.RandomValue),
		"Frees":         types.Gauge(s.Frees),
	}
}

func (s *Stats) GetCounters() map[string]types.Counter {
	return map[string]types.Counter{
		"PollCount": types.Counter(s.PollCount),
	}
}
