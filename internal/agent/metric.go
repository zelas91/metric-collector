package agent

import (
	"math/rand"
	"runtime"
)

type Stats struct {
	MemStats    *runtime.MemStats
	PollCount   int
	RandomValue int
}

func ReadStats(s *Stats) {
	runtime.ReadMemStats(s.MemStats)
	s.PollCount += 1
	s.RandomValue = rand.Int()
}
