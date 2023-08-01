package main

import "github.com/zelas91/metric-collector/internal/agent"

const (
	baseUrl = "http://localhost:8080/update"
)

func main() {
	agent.Run(2, 10, baseUrl)
}
