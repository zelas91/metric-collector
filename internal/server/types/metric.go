package types

type MetricTypeValue interface {
	isValue()
}
type Gauge float64

func (Gauge) isValue() {}

type Counter int64

func (Counter) isValue() {}

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)
