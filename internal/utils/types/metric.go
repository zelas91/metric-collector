package types

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)
