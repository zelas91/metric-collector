package main

import (
	"github.com/zelas91/metric-collector/internal/agent"
	"net/http"
	"sync"
	"time"
)

func main() {
	s := agent.NewStats()
	c := agent.HttpClient{
		Client: &http.Client{},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			<-time.After(2 * time.Second)
			s.ReadStats()

		}
	}()
	go func() {
		for {
			<-time.After(10 * time.Second)
			c.UpdateMetrics(s)

		}
	}()
	wg.Wait()
}
