package main

import (
	"flag"
	"fmt"
)

var addr *string
var pollInterval *int
var reportInterval *int

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	pollInterval = flag.Int("p", 2, " poll interval ")
	reportInterval = flag.Int("r", 10, " poll interval ")

}

type Config struct {
	BaseURL        string
	ReportInterval int
	PollInterval   int
}

func NewConfig() *Config {

	flag.Parse()
	baseURL := fmt.Sprintf("http://%s/update", *addr)
	return &Config{
		BaseURL:        baseURL,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
	}
}
