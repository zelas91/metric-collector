package types

type MetricType int

type Gauge struct {
	Value float64
}

type Counter struct {
	Value int64
}

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)
