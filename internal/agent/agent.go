package agent

import (
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
func (s *Stats) GetGauges() map[string]float64 {
	return map[string]float64{
		"Alloc":         float64(s.Alloc),
		"GCSys":         float64(s.GCSys),
		"HeapAlloc":     float64(s.HeapAlloc),
		"BuckHashSys":   float64(s.BuckHashSys),
		"GCCPUFraction": s.GCCPUFraction,
		"HeapIdle":      float64(s.HeapIdle),
		"HeapInuse":     float64(s.HeapInuse),
		"HeapObjects":   float64(s.HeapObjects),
		"HeapReleased":  float64(s.HeapReleased),
		"HeapSys":       float64(s.HeapSys),
		"LastGC":        float64(s.LastGC),
		"Lookups":       float64(s.Lookups),
		"MCacheInuse":   float64(s.MCacheInuse),
		"MCacheSys":     float64(s.MCacheSys),
		"MSpanInuse":    float64(s.MSpanInuse),
		"MSpanSys":      float64(s.MSpanSys),
		"Mallocs":       float64(s.Mallocs),
		"NextGC":        float64(s.NextGC),
		"NumForcedGC":   float64(s.NumForcedGC),
		"NumGC":         float64(s.NumGC),
		"OtherSys":      float64(s.OtherSys),
		"PauseTotalNs":  float64(s.PauseTotalNs),
		"StackInuse":    float64(s.StackInuse),
		"StackSys":      float64(s.StackSys),
		"Sys":           float64(s.Sys),
		"TotalAlloc":    float64(s.TotalAlloc),
		"RandomValue":   float64(s.RandomValue),
		"Frees":         float64(s.Frees),
	}
}

func (s *Stats) GetCounters() map[string]int64 {
	return map[string]int64{
		"PollCount": s.PollCount,
	}
}
