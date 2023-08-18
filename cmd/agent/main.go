package main

import (
	"github.com/zelas91/metric-collector/internal/agent"
)

func main() {
	conf := NewConfig()
	agent.Run(conf.PollInterval, conf.ReportInterval, conf.BaseURL)
}
