package repository

type MemRepository interface {
	AddMetric(name, typeMetric, value string)
	ReadMetric(name string) map[string]interface{}
}
