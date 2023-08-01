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
func (s *Stats) GetGauges() map[string]types.Gauge {
	return map[string]types.Gauge{
		"Alloc":         {Value: float64(s.Alloc)},
		"GCSys":         {Value: float64(s.GCSys)},
		"HeapAlloc":     {Value: float64(s.HeapAlloc)},
		"BuckHashSys":   {Value: float64(s.BuckHashSys)},
		"GCCPUFraction": {Value: s.GCCPUFraction},
		"HeapIdle":      {Value: float64(s.HeapIdle)},
		"HeapInuse":     {Value: float64(s.HeapInuse)},
		"HeapObjects":   {Value: float64(s.HeapObjects)},
		"HeapReleased":  {Value: float64(s.HeapReleased)},
		"HeapSys":       {Value: float64(s.HeapSys)},
		"LastGC":        {Value: float64(s.LastGC)},
		"Lookups":       {Value: float64(s.Lookups)},
		"MCacheInuse":   {Value: float64(s.MCacheInuse)},
		"MCacheSys":     {Value: float64(s.MCacheSys)},
		"MSpanInuse":    {Value: float64(s.MSpanInuse)},
		"MSpanSys":      {Value: float64(s.MSpanSys)},
		"Mallocs":       {Value: float64(s.Mallocs)},
		"NextGC":        {Value: float64(s.NextGC)},
		"NumForcedGC":   {Value: float64(s.NumForcedGC)},
		"NumGC":         {Value: float64(s.NumGC)},
		"OtherSys":      {Value: float64(s.OtherSys)},
		"PauseTotalNs":  {Value: float64(s.PauseTotalNs)},
		"StackInuse":    {Value: float64(s.StackInuse)},
		"StackSys":      {Value: float64(s.StackSys)},
		"Sys":           {Value: float64(s.Sys)},
		"TotalAlloc":    {Value: float64(s.TotalAlloc)},
		"RandomValue":   {Value: float64(s.RandomValue)},
	}
}

func (s *Stats) GetCounters() map[string]types.Counter {
	return map[string]types.Counter{
		"PoolCounter": {Value: s.PollCount},
	}
}
